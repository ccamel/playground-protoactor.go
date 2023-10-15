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
package core

import (
	"github.com/AsynkronIT/protoactor-go/actor"

	"github.com/ccamel/playground-protoactor.go/internal/system/core/sys"
	"github.com/ccamel/playground-protoactor.go/internal/system/core/usr"
)

type Actor struct{}

func (a *Actor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		_, _ = context.SpawnNamed(actor.PropsFromProducer(sys.New()), "sys")
		_, _ = context.SpawnNamed(actor.PropsFromProducer(usr.New()), "usr")
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	}
}

func New() actor.Producer {
	return func() actor.Actor {
		return &Actor{}
	}
}
