package booklend

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/ccamel/playground-protoactor.go/internal/app"
	"github.com/ccamel/playground-protoactor.go/internal/eventsourcing"
	"github.com/ccamel/playground-protoactor.go/internal/middleware"
)

type App struct {
	middleware.SpawnAwareMixin
}

func (a *App) Receive(context actor.Context) {
	if _, ok := context.Message().(*actor.Started); ok {
		a.SpawnNamedOrDie(context, eventsourcing.ManagerProps(actor.PropsFromProducer(New)), "book_lend")
	}
}

func init() {
	app.RegisterMust("book_lend_app", func(_ any) actor.Producer {
		return func() actor.Actor {
			return &App{}
		}
	})
}
