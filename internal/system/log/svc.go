// Copyright © 2020 Chris Camel <camel.christophe@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package log

import (
	"fmt"
	"io"
	"os"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/rs/zerolog"
)

type LoggerActor struct {
	out io.Writer
}

func (a *LoggerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	case *LogMessage:
		_, err := a.out.Write(msg.Message)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func New() actor.Producer {
	return func() actor.Actor {
		return &LoggerActor{
			out: zerolog.ConsoleWriter{
				Out: os.Stdout,
			},
		}
	}
}
