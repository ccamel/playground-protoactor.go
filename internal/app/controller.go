package app

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/ccamel/playground-protoactor.go/internal/middleware"
)

// Controller is the actor responsible for managing the lifecycle of the applications.
type Controller struct {
	middleware.SpawnAwareMixin
}

func (state *Controller) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		for k, v := range Seq() {
			state.SpawnNamedOrDie(context, actor.PropsFromProducer(v(nil)), k)
		}
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	}
}

func Props() *actor.Props {
	supervisor := actor.NewOneForOneStrategy(5, 1000, func(_ interface{}) actor.Directive {
		return actor.RestartDirective
	})

	return actor.
		PropsFromProducer(
			func() actor.Actor {
				return &Controller{}
			},
			actor.WithSupervisor(supervisor))
}
