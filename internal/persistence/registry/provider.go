package registry

import (
	"fmt"
	"net/url"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	persistence2 "github.com/ccamel/playground-protoactor.go/internal/persistence"
	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
	"github.com/ccamel/playground-protoactor.go/internal/util"
)

type provider struct {
	providerState persistence.ProviderState
}

func (p *provider) GetState() persistence.ProviderState {
	return p.providerState
}

type storeAdapter struct {
	store persistence2.Store
}

func (s *storeAdapter) GetSnapshotInterval() int {
	return s.store.GetSnapshotInterval()
}

func (s *storeAdapter) GetSnapshot(actorName string) (snapshot interface{}, eventIndex int, ok bool) {
	record, err := s.store.GetSnapshot(actorName)
	if err != nil {
		return nil, 0, false
	}

	message := &persistencev1.ConsiderSnapshot{
		Payload: record.Payload,
	}

	return message, int(record.Version), true //nolint: gosec // must adapt to the interface
}

func (s *storeAdapter) Restart() {
	s.store.Restart()
}

func (s *storeAdapter) PersistSnapshot(actorName string, eventIndex int, snapshot proto.Message) {
	payload, err := anypb.New(snapshot)
	if err != nil {
		log.Error().Err(err).Msg("Failed to persist snapshot")
	}

	entity := &persistencev1.SnapshotRecord{
		Id:               actorName,
		Type:             payload.TypeUrl,
		Version:          uint64(eventIndex), //nolint: gosec // must adapt to the interface
		StorageTimestamp: timestamppb.Now(),
		Payload:          payload,
	}

	s.store.PersistSnapshot(actorName, entity)
}

func (s *storeAdapter) DeleteSnapshots(actorName string, inclusiveToIndex int) {
	s.store.DeleteSnapshots(actorName, inclusiveToIndex)
}

func (s *storeAdapter) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e interface{})) {
	callbackAdapter := func(record *persistencev1.EventRecord) {
		message, err := record.Payload.UnmarshalNew()
		if err != nil {
			log.Error().
				Err(err).
				Str("id", record.Id).
				Uint64("index", record.Version).
				Msg("Failed to retrieve event")
		}
		callback(message)
	}
	s.store.GetEvents(actorName, eventIndexStart, eventIndexEnd, callbackAdapter)
}

func (s *storeAdapter) PersistEvent(actorName string, eventIndex int, event proto.Message) {
	payload, err := anypb.New(event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to persist event")
	}

	id := util.MakeULID()
	entity := &persistencev1.EventRecord{
		Id:               id.String(),
		Type:             payload.TypeUrl,
		StreamId:         actorName,
		Version:          uint64(eventIndex), //nolint: gosec // no overflow risk
		StorageTimestamp: timestamppb.Now(),
		Payload:          payload,
	}

	s.store.PersistEvent(actorName, entity)
}

func (s *storeAdapter) DeleteEvents(actorName string, inclusiveToIndex int) {
	s.store.DeleteEvents(actorName, inclusiveToIndex)
}

func NewProvider(system *actor.ActorSystem, uri persistence2.URI) (persistence.Provider, error) {
	if uri == "" {
		return nil, fmt.Errorf("persistence URI is required")
	}

	parsedURI, err := url.Parse(string(uri))
	if err != nil {
		return nil, err
	}

	factory, err := factories.GetFromURI(parsedURI)
	if err != nil {
		return nil, err
	}

	store, err := factory(system, parsedURI)
	if err != nil {
		return nil, err
	}

	return &provider{
		providerState: &storeAdapter{
			store: store,
		},
	}, nil
}
