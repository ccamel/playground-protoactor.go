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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogAware interface {
	SetLog(logger zerolog.Logger)
}

type LogAwareHolder struct {
	Log zerolog.Logger
}

func (state *LogAwareHolder) SetLog(logger zerolog.Logger) {
	state.Log = logger
}

type LogInjectorPlugin struct{}

func (p *LogInjectorPlugin) OnStart(ctx actor.ReceiverContext) {
	if p, ok := ctx.Actor().(LogAware); ok {
		p.SetLog(log.With().
			Str("pid", ctx.Self().GetId()).
			Logger())
	}
}

func (p *LogInjectorPlugin) OnOtherMessage(_ actor.ReceiverContext, _ *actor.MessageEnvelope) {}