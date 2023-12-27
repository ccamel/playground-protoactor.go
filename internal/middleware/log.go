package middleware

import (
	"fmt"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/rs/zerolog/log"
)

// LifecycleLogger is a middleware which logs lifecycle messages (started, stopped...).
func LifecycleLogger() actor.ReceiverMiddleware {
	return func(next actor.ReceiverFunc) actor.ReceiverFunc {
		return func(context actor.ReceiverContext, env *actor.MessageEnvelope) {
			accepted := func() bool {
				if _, ok := env.Message.(*actor.Started); ok {
					return true
				}

				if _, ok := env.Message.(*actor.Stopping); ok {
					return true
				}

				if _, ok := env.Message.(*actor.Stopped); ok {
					return true
				}

				if _, ok := env.Message.(*actor.Restarting); ok {
					return true
				}

				return false
			}()

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
