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
package persistence

import (
	"fmt"
	"strconv"

	"github.com/ccamel/playground-protoactor.go/internal/util"
	"github.com/golang/protobuf/descriptor" //nolint:staticcheck // use same version than protoactor library
	"github.com/golang/protobuf/proto"      //nolint:staticcheck // use same version than protoactor library
	"github.com/golang/protobuf/ptypes"     //nolint:staticcheck // use same version than protoactor library
	"github.com/rs/zerolog/log"
	bolt "go.etcd.io/bbolt"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

type BBoltProvider struct {
	snapshotInterval int
	db               *bolt.DB
}

func NewBBoltProvider(snapshotInterval int) (*Provider, error) {
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
		providerState: &BBoltProvider{
			snapshotInterval: snapshotInterval,
			db:               db,
		},
	}, nil
}

func (provider *BBoltProvider) Restart() {}

func (provider *BBoltProvider) GetSnapshotInterval() int {
	return provider.snapshotInterval
}

func (provider *BBoltProvider) GetSnapshot(actorName string) (interface{}, int, bool) {
	var message interface{}

	var eventIndex int

	err := provider.db.View(func(tx *bolt.Tx) error {
		buf := provider.
			snapshotsBucket(tx).
			Get([]byte(actorName))
		if buf == nil {
			return fmt.Errorf("snapshot %d not found: %w", eventIndex, ErrNotFound)
		}

		var entity Snapshot
		err := proto.Unmarshal(buf, &entity)
		if err != nil {
			return err
		}

		message = &ConsiderSnapshot{
			Payload: entity.Payload,
		}
		eventIndex = int(entity.Metadata.Index)

		return nil
	})

	return message, eventIndex, err == nil
}

func (provider *BBoltProvider) PersistSnapshot(actorName string, eventIndex int, snapshot proto.Message) {
	err := provider.db.Update(func(tx *bolt.Tx) error {
		payload, err := ptypes.MarshalAny(snapshot)
		if err != nil {
			return err
		}

		_, desc := descriptor.MessageDescriptorProto(snapshot)

		entity := &Snapshot{
			Metadata: &Snapshot_Metadata{
				Id:        actorName,
				Type:      *desc.Name,
				Timestamp: ptypes.TimestampNow(),
				Index:     uint64(eventIndex),
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

func (provider *BBoltProvider) GetEvents(actorName string, eventIndexStart int, callback func(e interface{})) {
	err := provider.db.View(func(tx *bolt.Tx) error {
		actorBucket := provider.
			eventsBucket(tx).
			Bucket([]byte(actorName))
		if actorBucket == nil {
			return nil
		}

		c := actorBucket.Cursor()

		for k, v := c.Seek(util.Itob(int64(eventIndexStart))); k != nil; k, v = c.Next() {
			var entity Event
			err := proto.Unmarshal(v, &entity)
			if err != nil {
				return err
			}

			var dynamic ptypes.DynamicAny
			err = ptypes.UnmarshalAny(entity.Payload, &dynamic)
			if err != nil {
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

func (provider *BBoltProvider) PersistEvent(actorName string, eventIndex int, event proto.Message) {
	err := provider.db.Update(func(tx *bolt.Tx) error {
		actorBucket, err := provider.
			eventsBucket(tx).
			CreateBucketIfNotExists([]byte(actorName))
		if err != nil {
			return err
		}

		payload, err := ptypes.MarshalAny(event)
		if err != nil {
			return err
		}

		id, _ := actorBucket.NextSequence()
		_, desc := descriptor.MessageDescriptorProto(event)

		entity := &Event{
			Metadata: &Event_Metadata{
				Id:        strconv.FormatUint(id, 10),
				Type:      *desc.Name,
				Timestamp: ptypes.TimestampNow(),
				Index:     uint64(eventIndex),
			},
			Payload: payload,
		}

		log.Info().Interface("entity", entity).Msg("Event saved")

		buf, err := proto.Marshal(entity)
		if err != nil {
			return err
		}

		err = actorBucket.Put(util.Itob(int64(eventIndex)), buf)

		return err
	})
	if err != nil { // TODO: use panic instead
		log.Error().Err(err).Msg("Failed to persist event")
	}
}

func (provider *BBoltProvider) eventsBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte("events"))
}

func (provider *BBoltProvider) snapshotsBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte("snapshots"))
}
