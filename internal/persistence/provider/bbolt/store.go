package bbolt

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/atomic"
	"google.golang.org/protobuf/proto"

	persistence "github.com/ccamel/playground-protoactor.go/internal/persistence"
	"github.com/ccamel/playground-protoactor.go/internal/persistence/stream"
	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
	"github.com/ccamel/playground-protoactor.go/internal/util"
)

var ErrNotFound = fmt.Errorf("not found")

type subscription struct {
	actor     *actor.PID
	predicate stream.EventPredicate
	handler   func(event proto.Message)
}

type Store struct {
	system           *actor.ActorSystem
	snapshotInterval int
	db               *bolt.DB
	muPublish        sync.Mutex
	subscribers      *sync.Map
}

var _ persistence.Store = (*Store)(nil)

func (s *Store) Restart() {}

func (s *Store) GetSnapshotInterval() int {
	return s.snapshotInterval
}

func (s *Store) GetSnapshot(actorName string) (*persistencev1.SnapshotRecord, error) {
	var record persistencev1.SnapshotRecord
	err := s.db.View(func(tx *bolt.Tx) error {
		buf := s.
			snapshotsBucket(tx).
			Get([]byte(actorName))
		if buf == nil {
			return fmt.Errorf("snapshot not found for actor %s: %w", actorName, ErrNotFound)
		}

		err := proto.Unmarshal(buf, &record)
		if err != nil {
			return err
		}

		return nil
	})

	return &record, err
}

func (s *Store) PersistSnapshot(actorName string, record *persistencev1.SnapshotRecord) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		buf, err := proto.Marshal(record)
		if err != nil {
			return err
		}

		err = s.
			snapshotsBucket(tx).
			Put([]byte(actorName), buf)

		return err
	})
	if err != nil { // TODO: use panic instead
		log.Error().Err(err).Msg("Failed to persist snapshot")
	}
}

//nolint:gosec // we need to make some dirt conversions to adapt to the interfaces
func (s *Store) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e *persistencev1.EventRecord)) {
	err := s.db.View(func(tx *bolt.Tx) error {
		actorBucket := s.
			eventsBucket(tx).
			Bucket([]byte(actorName))
		if actorBucket == nil {
			return nil
		}

		c := actorBucket.Cursor()

		for k, v := c.Seek(util.Itob(uint64(eventIndexStart))); k != nil &&
			(!(bytes.Compare(k, util.Itob(uint64(eventIndexEnd))) <= 0) || (eventIndexEnd == 0)); k, v = c.Next() {
			buf := s.eventsBucket(tx).Get(v)

			var record persistencev1.EventRecord
			if err := proto.Unmarshal(buf, &record); err != nil {
				return err
			}

			callback(&record)
		}

		return nil
	})
	if err != nil { // TODO: use panic instead
		log.Error().Err(err).Msg("Failed to retrieve events")
	}
}

func (s *Store) PersistEvent(actorName string, record *persistencev1.EventRecord) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		// store in the aggregate bucket the version number and the id of the record in the
		// events bucket.
		aggregateBucket, err := s.
			eventsBucket(tx).
			CreateBucketIfNotExists([]byte(actorName))
		if err != nil {
			return err
		}

		binID := []byte(record.Id)
		if err != nil {
			return err
		}

		record.SequenceNumber, _ = s.eventsBucket(tx).NextSequence()

		buf, err := proto.Marshal(record)
		if err != nil {
			return err
		}

		err = aggregateBucket.Put(util.Itob(record.Version), binID)
		if err != nil {
			return err
		}

		return s.
			eventsBucket(tx).
			Put(binID, buf)
	})
	if err != nil { // TODO: use panic instead
		log.Error().Err(err).Interface("record", record).Msg("Failed to persist event")

		return
	}

	s.publish(record)
}

func (s *Store) publish(event *persistencev1.EventRecord) {
	s.muPublish.Lock()
	defer s.muPublish.Unlock()

	s.subscribers.Range(func(_, value interface{}) bool {
		sub := value.(subscription)
		if sub.predicate(event) {
			sub.handler(event)
		}

		return true
	})
}

func (s *Store) Subscribe(pid *actor.PID, last *string, predicate stream.EventPredicate) stream.SubscriptionID {
	flag := atomic.NewBool(false)
	buffer := make([]interface{}, 0, 64)

	subscriptionID := uuid.NewString()
	s.subscribers.Store(
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
							s.system.Root.Send(pid, oldEvt)
						}
						buffer = nil
					}

					s.system.Root.Send(pid, event)
				}
			},
		},
	)

	go func() {
		defer func() {
			flag.Toggle()
		}()

		err := s.db.View(func(tx *bolt.Tx) error {
			c := s.eventsBucket(tx).Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				if last != nil && bytes.Compare(k, []byte(*last)) < 0 {
					continue
				}

				buf := s.eventsBucket(tx).Get(v)

				evt, err := unmarshallPayload(buf)
				if err != nil {
					return err
				}

				if predicate(evt.(*persistencev1.EventRecord)) {
					s.system.Root.Send(pid, evt)
				}
			}

			return nil
		})
		if err != nil {
			return
		}
	}()

	return stream.SubscriptionID(subscriptionID)
}

func (s *Store) Unsubscribe(_ stream.SubscriptionID) {
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) DeleteEvents(_ string, _ int) {
	// TODO: implement me!
}

func (s *Store) DeleteSnapshots(_ string, _ int) {
	// TODO: implement me!
}

// eventsBucket returns the bucket where all the events are stored in sequential order.
// In this bucket, a sub-bucket is created per aggregateId for quick retrieval of events for a considered aggregate.
func (s *Store) eventsBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte("events"))
}

func (s *Store) snapshotsBucket(tx *bolt.Tx) *bolt.Bucket {
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
