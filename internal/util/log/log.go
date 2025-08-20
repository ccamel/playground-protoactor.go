package log

import (
	"github.com/asynkron/protoactor-go/actor"
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
