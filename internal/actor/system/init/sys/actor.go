package sys

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/ccamel/playground-protoactor.go/internal/actor/system/log"
	"github.com/ccamel/playground-protoactor.go/internal/middleware"
)

type Actor struct {
	middleware.SpawnAwareMixin
}

func (a *Actor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		a.SpawnNamedOrDie(context, actor.PropsFromProducer(log.New()), "logger")
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
