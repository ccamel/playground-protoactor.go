// Copyright © 2020 Chris Camel <camel.christophe@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package middleware

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/persistence"
	"github.com/golang/protobuf/proto" //nolint:staticcheck // use same version than protoactor library
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
