package memory

import (
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
	"github.com/ccamel/playground-protoactor.go/internal/util"
)

var ErrNotFound = fmt.Errorf("not found")

type ProviderState struct {
	system           *actor.ActorSystem
	snapshotInterval int

	events    map[string][]*persistencev1.EventRecord
	snapshots map[string]*persistencev1.SnapshotRecord

	muEvent    sync.RWMutex
	muSnapshot sync.RWMutex

	subscribers *sync.Map
}

func (provider *ProviderState) Restart() {}

func (provider *ProviderState) GetSnapshotInterval() int {
	return provider.snapshotInterval
}

func (provider *ProviderState) GetSnapshot(actorName string) (interface{}, int, bool) {
	provider.muSnapshot.RLock()
	defer provider.muSnapshot.RUnlock()

	if record, ok := provider.snapshots[actorName]; ok {
		message := &persistencev1.ConsiderSnapshot{
			Payload: record.Payload,
		}
		eventIndex := int(record.Version)

		return message, eventIndex, true
	}

	return nil, 0, false
}

func (provider *ProviderState) PersistSnapshot(actorName string, eventIndex int, snapshot proto.Message) {
	payload, err := anypb.New(snapshot)
	if err != nil {
		log.
			Panic().
			Str("actor", actorName).
			Int("eventIndex", eventIndex).
			Err(err).
			Msg("Failed to create entity snapshot")
	}

	record := &persistencev1.SnapshotRecord{
		Id:               actorName,
		Type:             payload.TypeUrl,
		Version:          uint64(eventIndex),
		StorageTimestamp: timestamppb.Now(),
		Payload:          payload,
	}

	log.Info().Interface("record", record).Msg("Snapshot saved")

	provider.muSnapshot.Lock()
	provider.snapshots[actorName] = record
	provider.muSnapshot.Unlock()
}

func (provider *ProviderState) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e interface{})) {
	provider.muEvent.RLock()
	defer provider.muEvent.RUnlock()

	if events, ok := provider.events[actorName]; ok {
		for idx, event := range events {
			if eventIndexStart <= int(event.Version) && (eventIndexEnd == 0 || int(event.Version) <= eventIndexEnd) {
				payload, err := event.Payload.UnmarshalNew()
				if err != nil {
					log.
						Panic().
						Str("actor", actorName).
						Int("eventIndex", idx).
						Err(err).
						Msg("Failed to unmarshall entity event")
				}

				callback(payload)
			}
		}
	}
}

func (provider *ProviderState) PersistEvent(actorName string, eventIndex int, event proto.Message) {
	id := util.MakeULID()
	payload, err := anypb.New(event)
	if err != nil {
		log.
			Panic().
			Str("actor", actorName).
			Int("eventIndex", eventIndex).
			Err(err).
			Msg("Failed to create entity event")
	}

	record := &persistencev1.EventRecord{
		Id:               id.String(),
		Type:             payload.TypeUrl,
		StreamId:         actorName,
		Version:          uint64(eventIndex),
		SequenceNumber:   uint64(len(provider.events[actorName])),
		StorageTimestamp: timestamppb.Now(),
		Payload:          payload,
	}

	provider.events[actorName] = append(provider.events[actorName], record)

	log.Info().Interface("entity", record).Msg("Event saved")
}

func (provider *ProviderState) DeleteEvents(_ string, _ int) {
	// TODO: implement me!
}

func (provider *ProviderState) DeleteSnapshots(_ string, _ int) {
	// TODO: implement me!
}
