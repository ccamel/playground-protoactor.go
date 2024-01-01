package persistence

import (
	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
)

// Store is the interface that wraps the basic persistence operations.
type Store interface {
	GetSnapshot(actorName string) (snapshot *persistencev1.SnapshotRecord, err error)
	PersistSnapshot(actorName string, snapshot *persistencev1.SnapshotRecord)
	DeleteSnapshots(actorName string, inclusiveToIndex int)
	GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e *persistencev1.EventRecord))
	PersistEvent(actorName string, event *persistencev1.EventRecord)
	DeleteEvents(actorName string, inclusiveToIndex int)
	Restart()
	GetSnapshotInterval() int
}
