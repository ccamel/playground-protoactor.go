package init

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/ccamel/playground-protoactor.go/internal/actor/system/init/sys"
	"github.com/ccamel/playground-protoactor.go/internal/actor/system/init/usr"
)

type Actor struct{}

func (a *Actor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		_, _ = context.SpawnNamed(actor.PropsFromProducer(usr.New()), "usr")
		_, _ = context.SpawnNamed(actor.PropsFromProducer(sys.New()), "sys")
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
