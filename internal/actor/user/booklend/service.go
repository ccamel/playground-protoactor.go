package booklend

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"google.golang.org/genproto/googleapis/rpc/code"

	booklendv1 "github.com/ccamel/playground-protoactor.go/internal/actor/user/booklend/v1"
	"github.com/ccamel/playground-protoactor.go/internal/middleware"
)

type Service struct {
	middleware.LogAwareHolder
}

func (a *Service) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *booklendv1.RegisterBook:
		a.doCommand(context, msg.BookId)
	case *booklendv1.LendBook:
		a.doCommand(context, msg.BookId)
	case *booklendv1.ReturnBook:
		a.doCommand(context, msg.BookId)
	default:
		if context.Sender() != nil {
			context.Respond(&booklendv1.CommandStatus{
				Code:    code.Code_INVALID_ARGUMENT,
				Message: fmt.Sprintf("message %T is not supported", msg),
			})
		}
	}
}

// doCommand process the given command to the aggregate.
func (a *Service) doCommand(context actor.Context, id string) {
	book, err := getOrSpawn(context, id)
	if err != nil {
		context.Respond(&booklendv1.CommandStatus{
			Code:    code.Code_UNKNOWN,
			Message: err.Error(),
		})
	}

	context.Forward(book)
}

func getOrSpawn(context actor.Context, name string) (*actor.PID, error) {
	id := context.Self().Id + "/" + name
	for _, pid := range context.Children() {
		if pid.GetId() == id {
			return pid, nil
		}
	}

	return context.SpawnNamed(newAggregate(), name)
}

func NewService() *actor.Props {
	return actor.
		PropsFromProducer(func() actor.Actor {
			return &Service{}
		})
}
