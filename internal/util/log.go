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
package util

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rs/zerolog"
)

// LoggerFunc turns a function into an a zerolog marshaller.
type LoggerFunc func(e *zerolog.Event)

// MarshalZerologObject makes the LoggerFunc type a LogObjectMarshaler.
func (f LoggerFunc) MarshalZerologObject(e *zerolog.Event) {
	f(e)
}

// MapAsZerologObject converts a map into a LogObjectMarshaler.
func MapAsZerologObject(m map[string]string) LoggerFunc {
	return func(e *zerolog.Event) {
		for k, v := range m {
			e.Str(k, v)
		}
	}
}

// MessageEnvelopeAsZerologObject converts a map into a LogObjectMarshaler.
func MessageEnvelopeAsZerologObject(env *actor.MessageEnvelope) LoggerFunc {
	return LoggerFunc(func(e *zerolog.Event) {
		e.
			Object("", MapAsZerologObject(env.Header.ToMap()))
	})
}
