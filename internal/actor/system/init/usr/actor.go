package usr

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/ccamel/playground-protoactor.go/internal/app"
	"github.com/ccamel/playground-protoactor.go/internal/middleware"
)

type Actor struct {
	middleware.SpawnAwareMixin
}

func (state *Actor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		state.SpawnNamedOrDie(context, app.Props(), "app_mgr")
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
