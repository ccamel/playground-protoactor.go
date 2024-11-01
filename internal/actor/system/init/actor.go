package init

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/ccamel/playground-protoactor.go/internal/actor/system/init/sys"
	"github.com/ccamel/playground-protoactor.go/internal/actor/system/init/usr"
	"github.com/ccamel/playground-protoactor.go/internal/middleware"
)

type Actor struct {
	middleware.SpawnAwareMixin
}

func (a *Actor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		a.SpawnNamedOrDie(context, actor.PropsFromProducer(sys.New()), "sys")
		a.SpawnNamedOrDie(context, actor.PropsFromProducer(usr.New()), "usr")
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
