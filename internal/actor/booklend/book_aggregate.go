// Copyright © 2020 Chris Camel <camel.christophe@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package booklend

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/persistence"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/genproto/googleapis/rpc/code"

	booklendv1 "github.com/ccamel/playground-protoactor.go/internal/actor/booklend/v1"
	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
)

type BookAggregate struct {
	persistence.Mixin
	state *booklendv1.BookEntity
}

func (a *BookAggregate) Receive(context actor.Context) {
	a.handleMessage(context, context.Message())
}

//nolint:funlen // relax
func (a *BookAggregate) handleMessage(context actor.Context, message interface{}) {
	switch msg := message.(type) {
	case *actor.Started:
		a.state = &booklendv1.BookEntity{}
	case *actor.ReceiveTimeout:
		context.Stop(context.Self())
	case *persistence.RequestSnapshot:
		a.PersistSnapshot(a.state)
	case *persistencev1.ConsiderSnapshot:
		var dynamic ptypes.DynamicAny

		err := ptypes.UnmarshalAny(msg.Payload, &dynamic)
		if err != nil {
			panic(err)
		}

		a.state = dynamic.Message.(*booklendv1.BookEntity)
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
				Timestamp: ptypes.TimestampNow(),
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
				Timestamp:        ptypes.TimestampNow(),
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

		t2, err := ptypes.Timestamp(msg.Date)
		if err != nil {
			context.Respond(&booklendv1.CommandStatus{
				Code:    code.Code_UNKNOWN,
				Message: fmt.Sprintf("failed to parse date: %s", err.Error()),
			})

			break
		}

		t1, err := ptypes.Timestamp(a.state.Date)
		if err != nil {
			context.Respond(&booklendv1.CommandStatus{
				Code:    code.Code_UNKNOWN,
				Message: fmt.Sprintf("failed to parse date: %s", err.Error()),
			})

			break
		}

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
				Timestamp:    ptypes.TimestampNow(),
				By:           a.state.Borrower,
				Date:         msg.Date,
				LentDuration: ptypes.DurationProto(t2.Sub(t1)),
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

func (a *BookAggregate) applyAndReply(context actor.Context, response proto.Message, events ...proto.Message) {
	// save sender - issue https://github.com/AsynkronIT/protoactor-go/issues/256
	sender := context.Sender()

	for _, event := range events {
		a.handleMessage(context, event)
		a.PersistReceive(event)
	}

	if response != nil {
		context.Send(sender, response)
	}
}

func newBookAggregate() *actor.Props {
	return actor.
		PropsFromProducer(func() actor.Actor {
			return &BookAggregate{}
		})
}
