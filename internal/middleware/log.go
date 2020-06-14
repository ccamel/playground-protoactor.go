// Copyright © 2020 Chris Camel <camel.christophe@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package middleware

import (
	"fmt"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
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
				log.Info().
					Stringer("actor", context.Self()).
					Str("state", t).
					Msgf("actor <%s> %s", context.Self().Id, s)
			}

			next(context, env)
		}
	}
}
