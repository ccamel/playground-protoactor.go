package persistence

import (
	"github.com/asynkron/protoactor-go/actor"

	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
)

type EventPredicate func(event *persistencev1.EventRecord) bool

type SubscriptionID string

type EventStreamStore interface {
	// Subscribe does a subscription for events.
	Subscribe(pid *actor.PID, start *string, predicate EventPredicate) SubscriptionID
	Unsubscribe(SubscriptionID)
}
