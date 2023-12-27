package log

import (
	"fmt"
	"io"
	"os"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/rs/zerolog"

	logv1 "github.com/ccamel/playground-protoactor.go/internal/system/log/v1"
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
	case *logv1.LogMessage:
		_, err := a.out.Write(msg.Message)
		if err != nil {
			fmt.Println(err) //nolint:forbidigo // common pattern when using cobra library
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
