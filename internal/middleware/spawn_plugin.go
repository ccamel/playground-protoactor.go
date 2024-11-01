package middleware

import (
	"github.com/asynkron/protoactor-go/actor"
)

type spawnAware interface {
	LogAware
	SpawnNamedOrDie(ctx actor.Context, props *actor.Props, name string) *actor.PID
}

var _ spawnAware = (*SpawnAwareMixin)(nil)

type SpawnAwareMixin struct {
	LogAwareHolder
}

func (state *SpawnAwareMixin) SpawnNamedOrDie(ctx actor.Context, props *actor.Props, name string) *actor.PID {
	pid, err := ctx.SpawnNamed(props, name)
	if err != nil {
		state.Logger().Fatal().Err(err).Msgf("failed to spawn %s actor", name)
	}

	return pid
}

type SpawnInjectorPlugin struct {
	delegated LogInjectorPlugin
}

func (p *SpawnInjectorPlugin) OnStart(ctx actor.ReceiverContext) {
	p.delegated.OnStart(ctx)
}

func (p *SpawnInjectorPlugin) OnOtherMessage(ctx actor.ReceiverContext, msg *actor.MessageEnvelope) {
	p.delegated.OnOtherMessage(ctx, msg)
}
