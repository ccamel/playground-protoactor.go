package middleware

import (
	"github.com/asynkron/protoactor-go/actor"
)

func OptionalUsing(
	one actor.ReceiverMiddleware,
	predicate func(ctx actor.ReceiverContext, env *actor.MessageEnvelope) bool,
) func(next actor.ReceiverFunc) actor.ReceiverFunc {
	return func(two actor.ReceiverFunc) actor.ReceiverFunc {
		fn := func(ctx actor.ReceiverContext, env *actor.MessageEnvelope) {
			if predicate(ctx, env) {
				one(two)(ctx, env)
			} else {
				two(ctx, env)
			}
		}

		return fn
	}
}
