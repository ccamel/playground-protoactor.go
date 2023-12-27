package middleware

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"google.golang.org/protobuf/proto"
)

type persistent interface { // hack, as interface from protobuf is not public
	PersistReceive(message proto.Message)
}

// PersistenceUsing installs the persistence mixin only on Persistent actors.
// TODO: it would be preferable to have the peristence.Using function to be less strict regarding the
// nature of the actor.
func PersistenceUsing(provider persistence.Provider) func(next actor.ReceiverFunc) actor.ReceiverFunc {
	return OptionalUsing(
		persistence.Using(provider),
		func(ctx actor.ReceiverContext, env *actor.MessageEnvelope) bool {
			_, ok := ctx.Actor().(persistent)

			return ok
		},
	)
}
