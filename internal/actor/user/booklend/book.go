package booklend

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	booklendv1 "github.com/ccamel/playground-protoactor.go/internal/actor/user/booklend/v1"
	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
)

type Book struct {
	persistence.Mixin
	state *booklendv1.BookEntity
}

func (a *Book) Receive(context actor.Context) {
	a.handleMessage(context, context.Message())
}

//nolint:funlen // relax
func (a *Book) handleMessage(context actor.Context, message interface{}) {
	switch msg := message.(type) {
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
	case *booklendv1.RegisterBook:
		if a.state.Id != "" {
			context.Respond(&booklendv1.CommandStatus{
				Code:    code.Code_ALREADY_EXISTS,
				Message: fmt.Sprintf("book with id %s already exists.", msg.BookId),
			})

			break
		}

		a.applyAndReply(
			context,
			&booklendv1.CommandStatus{
				Code:    code.Code_OK,
				Message: fmt.Sprintf("book registered with id %s", msg.BookId),
			},
			&booklendv1.BookRegistered{
				Id:        msg.BookId,
				Timestamp: timestamppb.Now(),
				Title:     msg.Title,
				Isbn:      msg.Isbn,
			})
	case *booklendv1.LendBook:
		if msg.Borrower == "" {
			context.Respond(&booklendv1.CommandStatus{
				Code:    code.Code_INVALID_ARGUMENT,
				Message: fmt.Sprintf("command LendBook for book %s shall specify a borrower.", msg.BookId),
			})

			break
		}

		if a.state.Borrower != "" {
			context.Respond(&booklendv1.CommandStatus{
				Code:    code.Code_INVALID_ARGUMENT,
				Message: fmt.Sprintf("book with id %s is already lent.", msg.BookId),
			})

			break
		}

		a.applyAndReply(
			context,
			&booklendv1.CommandStatus{
				Code:    code.Code_OK,
				Message: fmt.Sprintf("book registered with id %s", msg.BookId),
			},
			&booklendv1.BookLent{
				Id:               msg.BookId,
				Timestamp:        timestamppb.Now(),
				Borrower:         msg.Borrower,
				Date:             msg.Date,
				ExpectedDuration: msg.ExpectedDuration,
			})
	case *booklendv1.ReturnBook:
		if a.state.Borrower == "" {
			context.Respond(&booklendv1.CommandStatus{
				Code:    code.Code_INVALID_ARGUMENT,
				Message: fmt.Sprintf("book with id %s has not been lent.", msg.BookId),
			})

			break
		}

		if !msg.Date.IsValid() {
			context.Respond(&booklendv1.CommandStatus{
				Code:    code.Code_UNKNOWN,
				Message: fmt.Sprintf("date %s is invalid", msg.Date.String()),
			})

			break
		}
		t2 := msg.Date.AsTime()

		if !a.state.Date.IsValid() {
			context.Respond(&booklendv1.CommandStatus{
				Code:    code.Code_UNKNOWN,
				Message: fmt.Sprintf("date %s is invalid", a.state.Date.String()),
			})

			break
		}
		t1 := a.state.Date.AsTime()

		if t2.Before(t1) {
			context.Respond(&booklendv1.CommandStatus{
				Code:    code.Code_INVALID_ARGUMENT,
				Message: fmt.Sprintf("book with id %s cannot be returned before being lent", msg.BookId),
			})

			break
		}

		a.applyAndReply(
			context,
			&booklendv1.CommandStatus{
				Code:    code.Code_OK,
				Message: fmt.Sprintf("book registered with id %s", msg.BookId),
			},
			&booklendv1.BookReturned{
				Id:           msg.BookId,
				Timestamp:    timestamppb.Now(),
				By:           a.state.Borrower,
				Date:         msg.Date,
				LentDuration: durationpb.New(t2.Sub(t1)),
			})

	case *booklendv1.BookRegistered:
		a.state.Id = msg.Id
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

func (a *Book) applyAndReply(context actor.Context, response proto.Message, events ...proto.Message) {
	// save sender - issue https://github.com/asynkron/protoactor-go/issues/256
	sender := context.Sender()

	for _, event := range events {
		a.handleMessage(context, event)
		a.PersistReceive(event)
	}

	if response != nil {
		context.Send(sender, response)
	}
}

func newAggregate() *actor.Props {
	return actor.
		PropsFromProducer(func() actor.Actor {
			return &Book{}
		})
}
