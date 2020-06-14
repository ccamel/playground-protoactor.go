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
	"github.com/AsynkronIT/protoactor-go/actor"
)

func OptionalUsing(one actor.ReceiverMiddleware, predicate func(ctx actor.ReceiverContext, env *actor.MessageEnvelope) bool) func(next actor.ReceiverFunc) actor.ReceiverFunc {
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
