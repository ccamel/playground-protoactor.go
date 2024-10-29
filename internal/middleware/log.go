package middleware

import (
	"fmt"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/rs/zerolog/log"
)

func isType[T any](typ any) bool {
	_, ok := typ.(*T)
	return ok
}

// LifecycleLogger is a middleware which logs lifecycle messages (started, stopped...).
func LifecycleLogger() actor.ReceiverMiddleware {
	return func(next actor.ReceiverFunc) actor.ReceiverFunc {
		return func(context actor.ReceiverContext, env *actor.MessageEnvelope) {
			accepted := isType[actor.Started](env.Message) ||
				isType[actor.Stopping](env.Message) ||
				isType[actor.Stopped](env.Message) ||
				isType[actor.Restarting](env.Message)

			if accepted {
				t := strings.TrimLeft(fmt.Sprintf("%T", env.Message), "*")
				s := strings.ToLower(t[strings.LastIndex(t, ".")+1:])

				logger := log.Info()
				if a, ok := context.Actor().(LogAware); ok {
					logger = a.Logger().Info()
				}

				logger.
					Str("actor", fmt.Sprintf("%s@%s", context.Self().Id, context.Self().Address)).
					Str("state", t).
					Msgf("actor <%s> %s", context.Self().Id, s)
			}

			next(context, env)
		}
	}
}
