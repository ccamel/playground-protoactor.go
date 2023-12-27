package middleware

import (
	"sync"

	"github.com/asynkron/protoactor-go/actor"
)

func SystemSync(wg *sync.WaitGroup) actor.ReceiverMiddleware {
	return func(next actor.ReceiverFunc) actor.ReceiverFunc {
		return func(context actor.ReceiverContext, env *actor.MessageEnvelope) {
			switch env.Message.(type) {
			case *actor.Started:
				wg.Add(1)
			case *actor.Stopped:
				wg.Done()
			}
			next(context, env)
		}
	}
}
