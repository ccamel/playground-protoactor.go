//nolint:gosec // we need to make some dirt conversions to adapt to the interfaces
package memory

import (
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/actor"

	persistence "github.com/ccamel/playground-protoactor.go/internal/persistence"
	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
)

type Store struct {
	system           *actor.ActorSystem
	snapshotInterval int

	events    map[string][]*persistencev1.EventRecord
	snapshots map[string]*persistencev1.SnapshotRecord

	muEvent    sync.RWMutex
	muSnapshot sync.RWMutex

	subscribers *sync.Map
}

var _ persistence.Store = (*Store)(nil)

func (s *Store) Restart() {}

func (s *Store) GetSnapshotInterval() int {
	return s.snapshotInterval
}

func (s *Store) GetSnapshot(actorName string) (*persistencev1.SnapshotRecord, error) {
	s.muSnapshot.RLock()
	defer s.muSnapshot.RUnlock()

	if record, ok := s.snapshots[actorName]; ok {
		return record, nil
	}

	return nil, fmt.Errorf("snapshot not found for actor %s", actorName)
}

func (s *Store) PersistSnapshot(actorName string, record *persistencev1.SnapshotRecord) {
	s.muSnapshot.Lock()
	s.snapshots[actorName] = record
	s.muSnapshot.Unlock()
}

func (s *Store) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e *persistencev1.EventRecord)) {
	s.muEvent.RLock()
	defer s.muEvent.RUnlock()

	if events, ok := s.events[actorName]; ok {
		for _, event := range events {
			if eventIndexStart <= int(event.Version) && (eventIndexEnd == 0 || int(event.Version) <= eventIndexEnd) {
				callback(event)
			}
		}
	}
}

func (s *Store) PersistEvent(actorName string, record *persistencev1.EventRecord) {
	s.events[actorName] = append(s.events[actorName], record)
}

func (s *Store) DeleteEvents(_ string, _ int) {
	// TODO: implement me!
}

func (s *Store) DeleteSnapshots(_ string, _ int) {
	// TODO: implement me!
}
