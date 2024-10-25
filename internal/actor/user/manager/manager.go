package manager

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"google.golang.org/genproto/googleapis/rpc/code"

	eventsourcingv1 "github.com/ccamel/playground-protoactor.go/internal/eventsourcing/v1"
	"github.com/ccamel/playground-protoactor.go/internal/middleware"
)

type Manager struct {
	middleware.LogAwareHolder

	entityProps *actor.Props
}

func (a *Manager) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case eventsourcingv1.Command:
		a.doCommand(context, msg)
		return
	case actor.SystemMessage, actor.AutoReceiveMessage:
		// ignore
	default:
		if context.Sender() != nil {
			context.Respond(&eventsourcingv1.CommandStatus{
				Code:    code.Code_INVALID_ARGUMENT,
				Message: fmt.Sprintf("message %T is not supported", context.Message()),
			})
		}
	}
}

// doCommand process the given command to the aggregate.
func (a *Manager) doCommand(context actor.Context, cmd eventsourcingv1.Command) {
	entity, err := a.getOrSpawn(context, cmd.GetBase().AggregateId)
	if err != nil {
		context.Respond(&eventsourcingv1.CommandStatus{
			Code:    code.Code_UNKNOWN,
			Message: err.Error(),
		})
		return
	}

	context.Forward(entity)
}

func (a *Manager) getOrSpawn(context actor.Context, name string) (*actor.PID, error) {
	id := context.Self().Id + "/" + name
	for _, pid := range context.Children() {
		if pid.GetId() == id {
			return pid, nil
		}
	}

	return context.SpawnNamed(a.entityProps, name)
}

func Props(entityProps *actor.Props) *actor.Props {
	supervisor := actor.NewOneForOneStrategy(10, 1000, func(_ interface{}) actor.Directive {
		return actor.RestartDirective
	})

	return actor.
		PropsFromProducer(
			func() actor.Actor {
				return &Manager{
					entityProps: entityProps,
				}
			},
			actor.WithSupervisor(supervisor))
}
