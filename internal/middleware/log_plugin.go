package middleware

import (
	"bytes"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/rs/zerolog"

	logv1 "github.com/ccamel/playground-protoactor.go/internal/actor/system/log/v1"
)

type LogAware interface {
	Logger() *zerolog.Logger
	SetLog(logger zerolog.Logger)
}

type LogAwareHolder struct {
	Log zerolog.Logger
}

func (state *LogAwareHolder) SetLog(logger zerolog.Logger) {
	state.Log = logger
}

func (state *LogAwareHolder) Logger() *zerolog.Logger {
	return &state.Log
}

type LogInjectorPlugin struct{}

func (p *LogInjectorPlugin) OnStart(ctx actor.ReceiverContext) {
	if p, ok := ctx.Actor().(LogAware); ok {
		p.SetLog(
			zerolog.
				New(&loggerActor{
					root: ctx.ActorSystem(),
				}).
				With().
				Str("pid", ctx.Self().GetId()).
				Timestamp().
				Logger())
	}
}

func (p *LogInjectorPlugin) OnOtherMessage(_ actor.ReceiverContext, _ *actor.MessageEnvelope) {}

type loggerActor struct {
	root *actor.ActorSystem
	buf  bytes.Buffer
}

func (l *loggerActor) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if b != '\n' {
			l.buf.WriteByte(b)
		} else {
			pid := l.root.NewLocalPID("init/sys/logger")
			l.root.Root.Send(pid, &logv1.LogMessage{Message: l.buf.Bytes()})
			l.buf.Reset()
		}
	}

	return len(p), nil
}
