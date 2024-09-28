package system

import (
	SYS "syscall"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/actor/middleware/propagator"
	"github.com/asynkron/protoactor-go/plugin"
	"github.com/rs/zerolog/log"
	DEATH "github.com/vrecan/death"

	// Register bbolt persistence provider.
	_ "github.com/ccamel/playground-protoactor.go/internal/persistence/provider/bbolt"
	// Register memory persistence provider.
	_ "github.com/ccamel/playground-protoactor.go/internal/persistence/provider/memory"

	i "github.com/ccamel/playground-protoactor.go/internal/actor/system/init"
	"github.com/ccamel/playground-protoactor.go/internal/middleware"
	"github.com/ccamel/playground-protoactor.go/internal/persistence/registry"
)

const (
	// passivationTimeout is the duration after which an actor is passivated if it has not received any message.
	passivationTimeout = 5 * time.Second
)

type System struct {
	rootContext *actor.RootContext
	initPid     *actor.PID
}

func (s System) InitContext() *actor.RootContext {
	return s.rootContext
}

func (s System) Wait() {
	death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM)
	death.WaitForDeathWithFunc(func() {
		log.Info().
			Str("actor", "/").
			Msgf("ctrl-c received, stopping actor <%s>", s.initPid.String())
		err := s.rootContext.StopFuture(s.initPid).Wait()
		if err != nil {
			log.Error().
				Str("actor", "/").
				Err(err).
				Msgf("error while waiting for system shutdown")
		}
	})
}

func Boot(config Config) (*System, error) {
	log.Info().
		Str("actor", "/").
		Msg("booting the system...")

	system := actor.NewActorSystem()

	log.Info().
		Str("actor", "/").
		Str("registryAddress", system.ProcessRegistry.Address).
		Msg("system started")

	provider, err := registry.NewProvider(system, config.PersistenceURI)
	if err != nil {
		return nil, err
	}

	rootContext := system.
		Root.
		WithGuardian(actor.RestartingSupervisorStrategy()).
		WithSpawnMiddleware(
			propagator.New().
				WithItselfForwarded().
				WithReceiverMiddleware(
					plugin.Use(&middleware.LogInjectorPlugin{}),
					plugin.Use(&plugin.PassivationPlugin{Duration: passivationTimeout}),
					middleware.LifecycleLogger(),
					middleware.PersistenceUsing(provider),
				).
				SpawnMiddleware)

	props := actor.PropsFromProducer(i.New()).
		Configure(
			actor.WithSupervisor(actor.RestartingSupervisorStrategy()),
		)

	pid, err := rootContext.SpawnNamed(props, "init")
	if err != nil {
		return nil, err
	}

	return &System{
		rootContext: rootContext,
		initPid:     pid,
	}, nil
}
