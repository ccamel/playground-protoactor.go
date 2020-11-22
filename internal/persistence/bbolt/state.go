// Copyright © 2020 Chris Camel <camel.christophe@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package bbolt

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	p "github.com/AsynkronIT/protoactor-go/persistence"
	"github.com/ccamel/playground-protoactor.go/internal/persistence"
	"github.com/ccamel/playground-protoactor.go/internal/util"
	"github.com/golang/protobuf/proto"  //nolint:staticcheck // use same version than protoactor library
	"github.com/golang/protobuf/ptypes" //nolint:staticcheck // use same version than protoactor library
	"github.com/oklog/ulid"
	"github.com/rs/zerolog/log"
	bolt "go.etcd.io/bbolt"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

type ProviderState struct {
	snapshotInterval int
	db               *bolt.DB
	mu               sync.Mutex
	entropy          io.Reader
}

func NewProvider(snapshotInterval int) (p.Provider, error) {
	db, err := bolt.Open("./my-db", 0666, nil)
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
			snapshotInterval: snapshotInterval,
			db:               db,
			entropy:          ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0), //nolint:gosec
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

		var entity persistence.Snapshot
		err := proto.Unmarshal(buf, &entity)
		if err != nil {
			return err
		}

		message = &persistence.ConsiderSnapshot{
			Payload: entity.Payload,
		}
		eventIndex = int(entity.Metadata.Version)

		return nil
	})

	return message, eventIndex, err == nil
}

func (provider *ProviderState) PersistSnapshot(actorName string, eventIndex int, snapshot proto.Message) {
	err := provider.db.Update(func(tx *bolt.Tx) error {
		payload, err := ptypes.MarshalAny(snapshot)
		if err != nil {
			return err
		}

		entity := &persistence.Snapshot{
			Id: actorName,
			Metadata: &persistence.Snapshot_Metadata{
				StorageTimestamp: ptypes.TimestampNow(),
				Version:          uint64(eventIndex),
			},
			Payload: payload,
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

func (provider *ProviderState) GetEvents(actorName string, eventIndexStart int, callback func(e interface{})) {
	err := provider.db.View(func(tx *bolt.Tx) error {
		actorBucket := provider.
			eventsBucket(tx).
			Bucket([]byte(actorName))
		if actorBucket == nil {
			return nil
		}

		c := actorBucket.Cursor()

		for k, v := c.Seek(util.Itob(int64(eventIndexStart))); k != nil; k, v = c.Next() {
			buf := provider.eventsBucket(tx).Get(v)

			var entity persistence.Event
			if err := proto.Unmarshal(buf, &entity); err != nil {
				return err
			}

			var dynamic ptypes.DynamicAny
			if err := ptypes.UnmarshalAny(entity.Payload, &dynamic); err != nil {
				return err
			}

			callback(dynamic.Message)
		}

		return nil
	})
	if err != nil { // TODO: use panic instead
		log.Error().Err(err).Msg("Failed to retrieve events")
	}
}

func (provider *ProviderState) PersistEvent(actorName string, eventIndex int, event proto.Message) {
	err := provider.db.Update(func(tx *bolt.Tx) error {
		provider.mu.Lock()
		id := ulid.MustNew(ulid.Timestamp(time.Now()), provider.entropy)
		provider.mu.Unlock()

		binID, err := id.MarshalBinary()
		if err != nil {
			return err
		}

		// store in the aggregate bucket the version number and the id of the record in the
		// events bucket.
		aggregateBucket, err := provider.
			eventsBucket(tx).
			CreateBucketIfNotExists([]byte(actorName))
		if err != nil {
			return err
		}

		err = aggregateBucket.Put(util.Itob(int64(eventIndex)), binID)
		if err != nil {
			return err
		}

		// store in the events bucket the event
		payload, err := ptypes.MarshalAny(event)
		if err != nil {
			return err
		}

		entity := &persistence.Event{
			Id: id.String(),
			Metadata: &persistence.Event_Metadata{
				StorageTimestamp: ptypes.TimestampNow(),
				Version:          uint64(eventIndex),
			},
			Payload: payload,
		}

		log.Info().Interface("entity", entity).Msg("Event saved")

		buf, err := proto.Marshal(entity)
		if err != nil {
			return err
		}

		err = provider.
			eventsBucket(tx).
			Put(binID, buf)

		return err
	})
	if err != nil { // TODO: use panic instead
		log.Error().Err(err).Msg("Failed to persist event")
	}
}

// eventsBucket returns the bucket where all the events are stored in sequential order.
// In this bucket, a sub-bucket is created per aggregateId for quick retrieval of events for a considered aggregate.
func (provider *ProviderState) eventsBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte("events"))
}

func (provider *ProviderState) snapshotsBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte("snapshots"))
}
