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
	booklendv1 "github.com/ccamel/playground-protoactor.go/internal/actor/booklend/v1"
	"google.golang.org/genproto/googleapis/rpc/code"

	"github.com/ccamel/playground-protoactor.go/internal/middleware"
)

type BookService struct {
	middleware.LogAwareHolder
}

func (a *BookService) Receive(context actor.Context) {
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
func (a *BookService) doCommand(context actor.Context, id string) {
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

	return context.SpawnNamed(newBookAggregate(), name)
}

func NewBookCommandHandler() *actor.Props {
	return actor.
		PropsFromProducer(func() actor.Actor {
			return &BookService{}
		})
}
