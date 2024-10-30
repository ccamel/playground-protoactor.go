package booklend

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/asynkron/protoactor-go/plugin"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	booklendv1 "github.com/ccamel/playground-protoactor.go/internal/actor/user/booklend/v1"
	eventsourcingv1 "github.com/ccamel/playground-protoactor.go/internal/eventsourcing/v1"
	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
)

type Book struct {
	persistence.Mixin
	plugin.PassivationHolder
	state *booklendv1.BookEntity
}

func (a *Book) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		a.state = &booklendv1.BookEntity{}
	case *actor.ReceiveTimeout:
		context.Stop(context.Self())
	case *persistence.RequestSnapshot:
		a.PersistSnapshot(a.state)
	case *persistencev1.ConsiderSnapshot:
		entity := new(booklendv1.BookEntity)
		if err := msg.Payload.UnmarshalTo(entity); err != nil {
			panic(err)
		}

		a.state = entity
	case eventsourcingv1.Command:
		status, event := a.handleCommand(msg)
		// save sender - issue https://github.com/asynkron/protoactor-go/issues/256
		sender := context.Sender()
		a.PersistReceive(event.(proto.Message))
		a.handleEvent(event)

		if status != nil {
			context.Send(sender, status)
		}
	case eventsourcingv1.Event:
		a.handleEvent(msg)
	}
}

//nolint:funlen
func (a *Book) handleCommand(cmd eventsourcingv1.Command) (*eventsourcingv1.CommandStatus, eventsourcingv1.Event) {
	switch cmd := cmd.(type) {
	case *booklendv1.RegisterBook:
		if a.state.Id != "" {
			return &eventsourcingv1.CommandStatus{
				Code:    code.Code_ALREADY_EXISTS,
				Message: fmt.Sprintf("book with id %s already exists.", cmd.Base.AggregateId),
			}, nil
		}

		return &eventsourcingv1.CommandStatus{
				Code:    code.Code_OK,
				Message: fmt.Sprintf("book registered with id %s", cmd.Base.AggregateId),
			}, &booklendv1.BookRegistered{
				Base:      &eventsourcingv1.EventBase{Id: cmd.Base.AggregateId},
				Timestamp: timestamppb.Now(),
				Title:     cmd.Title,
				Isbn:      cmd.Isbn,
			}
	case *booklendv1.LendBook:
		if cmd.Borrower == "" {
			return &eventsourcingv1.CommandStatus{
				Code:    code.Code_INVALID_ARGUMENT,
				Message: fmt.Sprintf("command LendBook for book %s shall specify a borrower.", cmd.Base.AggregateId),
			}, nil
		}

		if a.state.Borrower != "" {
			return &eventsourcingv1.CommandStatus{
				Code:    code.Code_INVALID_ARGUMENT,
				Message: fmt.Sprintf("book with id %s is already lent.", cmd.Base.AggregateId),
			}, nil
		}

		return &eventsourcingv1.CommandStatus{
				Code:    code.Code_OK,
				Message: fmt.Sprintf("book lent with id %s", cmd.Base.AggregateId),
			}, &booklendv1.BookLent{
				Base:             &eventsourcingv1.EventBase{Id: cmd.Base.AggregateId},
				Timestamp:        timestamppb.Now(),
				Borrower:         cmd.Borrower,
				Date:             cmd.Date,
				ExpectedDuration: cmd.ExpectedDuration,
			}
	case *booklendv1.ReturnBook:
		if a.state.Borrower == "" {
			return &eventsourcingv1.CommandStatus{
				Code:    code.Code_INVALID_ARGUMENT,
				Message: fmt.Sprintf("book with id %s has not been lent.", cmd.Base.AggregateId),
			}, nil
		}

		if !cmd.Date.IsValid() {
			return &eventsourcingv1.CommandStatus{
				Code:    code.Code_UNKNOWN,
				Message: fmt.Sprintf("date %s is invalid", cmd.Date.String()),
			}, nil
		}
		t2 := cmd.Date.AsTime()

		if !a.state.Date.IsValid() {
			return &eventsourcingv1.CommandStatus{
				Code:    code.Code_UNKNOWN,
				Message: fmt.Sprintf("date %s is invalid", a.state.Date.String()),
			}, nil
		}
		t1 := a.state.Date.AsTime()

		if t2.Before(t1) {
			return &eventsourcingv1.CommandStatus{
				Code:    code.Code_INVALID_ARGUMENT,
				Message: fmt.Sprintf("book with id %s cannot be returned before being lent", cmd.Base.AggregateId),
			}, nil
		}

		return &eventsourcingv1.CommandStatus{
				Code:    code.Code_OK,
				Message: fmt.Sprintf("book returned with id %s", cmd.Base.AggregateId),
			}, &booklendv1.BookReturned{
				Base:         &eventsourcingv1.EventBase{Id: cmd.Base.AggregateId},
				Timestamp:    timestamppb.Now(),
				By:           a.state.Borrower,
				Date:         cmd.Date,
				LentDuration: durationpb.New(t2.Sub(t1)),
			}
	}

	return &eventsourcingv1.CommandStatus{
		Code:    code.Code_INVALID_ARGUMENT,
		Message: fmt.Sprintf("unsupported command %T received", cmd),
	}, nil
}

func (a *Book) handleEvent(msg eventsourcingv1.Event) {
	switch msg := msg.(type) {
	case *booklendv1.BookRegistered:
		a.state.Id = msg.Base.Id
		a.state.Isbn = msg.Isbn
		a.state.Title = msg.Title

	case *booklendv1.BookLent:
		a.state.Borrower = msg.Borrower
		a.state.Date = msg.Date
		a.state.ExpectedDuration = msg.ExpectedDuration

	case *booklendv1.BookReturned:
		a.state.Borrower = ""
	}
}

func New() actor.Actor {
	return &Book{}
}
