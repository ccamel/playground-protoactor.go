package bbolt

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	p "github.com/asynkron/protoactor-go/persistence"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/atomic"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ccamel/playground-protoactor.go/internal/persistence"
	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
	"github.com/ccamel/playground-protoactor.go/internal/util"
)

var ErrNotFound = fmt.Errorf("not found")

type subscription struct {
	actor     *actor.PID
	predicate persistence.EventPredicate
	handler   func(event proto.Message)
}

type ProviderState struct {
	system           *actor.ActorSystem
	snapshotInterval int
	db               *bolt.DB
	muID             sync.Mutex
	muPublish        sync.Mutex
	entropy          io.Reader
	subscribers      *sync.Map
}

func NewProvider(system *actor.ActorSystem, path string, snapshotInterval int) (p.Provider, error) {
	db, err := bolt.Open(path, 0o666, nil)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("bd", db.Path()).
		Msg("event store opened")

	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("events")); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("snapshots")); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Provider{
		providerState: &ProviderState{
			system:           system,
			snapshotInterval: snapshotInterval,
			db:               db,
			entropy:          ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0), //nolint:gosec
			subscribers:      &sync.Map{},
		},
	}, nil
}

func (provider *ProviderState) Restart() {}

func (provider *ProviderState) GetSnapshotInterval() int {
	return provider.snapshotInterval
}

func (provider *ProviderState) GetSnapshot(actorName string) (interface{}, int, bool) {
	var message interface{}

	var eventIndex int

	err := provider.db.View(func(tx *bolt.Tx) error {
		buf := provider.
			snapshotsBucket(tx).
			Get([]byte(actorName))
		if buf == nil {
			return fmt.Errorf("snapshot %d not found: %w", eventIndex, ErrNotFound)
		}

		var record persistencev1.SnapshotRecord
		err := proto.Unmarshal(buf, &record)
		if err != nil {
			return err
		}

		message = &persistencev1.ConsiderSnapshot{
			Payload: record.Payload,
		}
		eventIndex = int(record.Version)

		return nil
	})

	return message, eventIndex, err == nil
}

func (provider *ProviderState) PersistSnapshot(actorName string, eventIndex int, snapshot proto.Message) {
	err := provider.db.Update(func(tx *bolt.Tx) error {
		payload, err := anypb.New(snapshot)
		if err != nil {
			return err
		}

		entity := &persistencev1.SnapshotRecord{
			Id:               actorName,
			Type:             payload.TypeUrl,
			Version:          uint64(eventIndex),
			StorageTimestamp: timestamppb.Now(),
			Payload:          payload,
		}

		log.Info().Interface("entity", entity).Msg("Snapshot saved")

		buf, err := proto.Marshal(entity)
		if err != nil {
			return err
		}

		err = provider.
			snapshotsBucket(tx).
			Put([]byte(actorName), buf)

		return err
	})
	if err != nil { // TODO: use panic instead
		log.Error().Err(err).Msg("Failed to persist snapshot")
	}
}

func (provider *ProviderState) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e interface{})) {
	err := provider.db.View(func(tx *bolt.Tx) error {
		actorBucket := provider.
			eventsBucket(tx).
			Bucket([]byte(actorName))
		if actorBucket == nil {
			return nil
		}

		c := actorBucket.Cursor()

		for k, v := c.Seek(util.Itob(int64(eventIndexStart))); k != nil &&
			(!(bytes.Compare(k, util.Itob(int64(eventIndexEnd))) <= 0) || (eventIndexEnd == 0)); k, v = c.Next() {
			buf := provider.eventsBucket(tx).Get(v)

			i, err := unmarshallPayload(buf)
			if err != nil {
				return err
			}

			callback(i)
		}

		return nil
	})
	if err != nil { // TODO: use panic instead
		log.Error().Err(err).Msg("Failed to retrieve events")
	}
}

func (provider *ProviderState) PersistEvent(actorName string, eventIndex int, event proto.Message) {
	id, entity, err := func() (ulid.ULID, *persistencev1.EventRecord, error) {
		provider.muID.Lock()
		id, err := ulid.New(ulid.Timestamp(time.Now()), provider.entropy)
		provider.muID.Unlock()

		if err != nil {
			return ulid.ULID{}, nil, err
		}

		payload, err := anypb.New(event)
		if err != nil {
			return id, nil, err
		}

		return id, &persistencev1.EventRecord{
			Id:               id.String(),
			Type:             payload.TypeUrl,
			StreamId:         actorName,
			Version:          uint64(eventIndex),
			StorageTimestamp: timestamppb.Now(),
			Payload:          payload,
		}, nil
	}()
	if err != nil { // TODO: use panic instead
		log.Error().Err(err).Msg("Failed to create entity event")

		return
	}

	err = provider.db.Update(func(tx *bolt.Tx) error {
		// store in the aggregate bucket the version number and the id of the record in the
		// events bucket.
		aggregateBucket, err := provider.
			eventsBucket(tx).
			CreateBucketIfNotExists([]byte(actorName))
		if err != nil {
			return err
		}

		binID, err := id.MarshalBinary()
		if err != nil {
			return err
		}

		entity.SequenceNumber, _ = provider.eventsBucket(tx).NextSequence()

		buf, err := proto.Marshal(entity)
		if err != nil {
			return err
		}

		err = aggregateBucket.Put(util.Itob(int64(eventIndex)), binID)
		if err != nil {
			return err
		}

		return provider.
			eventsBucket(tx).
			Put(binID, buf)
	})

	if err != nil { // TODO: use panic instead
		log.Error().Err(err).Interface("entity", entity).Msg("Failed to persist event")

		return
	}

	provider.publish(entity)

	log.Info().Interface("entity", entity).Msg("Event saved")
}

func (provider *ProviderState) publish(event *persistencev1.EventRecord) {
	provider.muPublish.Lock()
	defer provider.muPublish.Unlock()

	provider.subscribers.Range(func(_, value interface{}) bool {
		sub := value.(subscription)
		if sub.predicate(event) {
			sub.handler(event)
		}

		return true
	})
}

func (provider *ProviderState) Subscribe(pid *actor.PID, last *string, predicate persistence.EventPredicate) persistence.SubscriptionID {
	flag := atomic.NewBool(false)
	buffer := make([]interface{}, 0, 64)

	subscriptionID := uuid.NewString()
	provider.subscribers.Store(
		subscriptionID,
		subscription{
			actor:     pid,
			predicate: predicate,
			handler: func(event proto.Message) {
				switch flag.Load() {
				case false:
					buffer = append(buffer, event)
				case true:
					if len(buffer) != 0 {
						for _, oldEvt := range buffer {
							provider.system.Root.Send(pid, oldEvt)
						}
						buffer = nil
					}

					provider.system.Root.Send(pid, event)
				}
			},
		},
	)

	go func() {
		defer func() {
			flag.Toggle()
		}()

		err := provider.db.View(func(tx *bolt.Tx) error {
			c := provider.eventsBucket(tx).Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				if last != nil && bytes.Compare(k, []byte(*last)) < 0 {
					continue
				}

				buf := provider.eventsBucket(tx).Get(v)

				evt, err := unmarshallPayload(buf)
				if err != nil {
					return err
				}

				if predicate(evt.(*persistencev1.EventRecord)) {
					provider.system.Root.Send(pid, evt)
				}
			}

			return nil
		})
		if err != nil {
			return
		}
	}()

	return persistence.SubscriptionID(subscriptionID)
}

func (provider *ProviderState) Unsubscribe(_ persistence.SubscriptionID) {
}

func (provider *ProviderState) Close() error {
	return provider.db.Close()
}

func (provider *ProviderState) DeleteEvents(_ string, _ int) {
	// TODO: implement me!
}

func (provider *ProviderState) DeleteSnapshots(_ string, _ int) {
	// TODO: implement me!
}

// eventsBucket returns the bucket where all the events are stored in sequential order.
// In this bucket, a sub-bucket is created per aggregateId for quick retrieval of events for a considered aggregate.
func (provider *ProviderState) eventsBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte("events"))
}

func (provider *ProviderState) snapshotsBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte("snapshots"))
}

func unmarshallPayload(buf []byte) (interface{}, error) {
	var entity persistencev1.EventRecord
	if err := proto.Unmarshal(buf, &entity); err != nil {
		return nil, err
	}

	message, err := entity.Payload.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	return message, nil
}
