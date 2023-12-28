package system

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	SYS "syscall"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/actor/middleware/propagator"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/asynkron/protoactor-go/plugin"
	"github.com/rs/zerolog/log"
	DEATH "github.com/vrecan/death"

	"github.com/ccamel/playground-protoactor.go/internal/middleware"
	"github.com/ccamel/playground-protoactor.go/internal/persistence/bbolt"
	"github.com/ccamel/playground-protoactor.go/internal/system/core"
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

	provider, err := getPersistenceProvider(system, config.PersistenceURI)
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
					middleware.LifecycleLogger(),
					middleware.PersistenceUsing(provider),
				).
				SpawnMiddleware)

	props := actor.PropsFromProducer(core.New()).
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

func getPersistenceProvider(system *actor.ActorSystem, uri URI) (persistence.Provider, error) {
	if uri == "" {
		return nil, fmt.Errorf("persistence URI is required")
	}

	parsedURI, err := url.Parse(string(uri))
	if err != nil {
		return nil, err
	}

	switch parsedURI.Scheme {
	case "db":
		parts := strings.Split(parsedURI.Opaque, ":")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid persistence URI: %s", uri)
		}
		head := parts[0]
		tail := parts[1]

		switch head {
		case "bbolt":
			path, err := url.PathUnescape(tail)
			if err != nil {
				return nil, fmt.Errorf("invalid persistence URI: %s. %w", uri, err)
			}

			snapshotInterval := 3
			if snapshotIntervalStr := parsedURI.Query().Get("snapshotInterval"); snapshotIntervalStr != "" {
				if snapshotInterval, err = strconv.Atoi(snapshotIntervalStr); err != nil {
					return nil, fmt.Errorf("invalid snapshotInterval value: %s. %w", snapshotIntervalStr, err)
				}
			}

			return bbolt.NewProvider(system, path, snapshotInterval)
		default:
			return nil, fmt.Errorf("unsupported database: %s", head)
		}
	default:
		return nil, fmt.Errorf("invalid persistence URI: %s", uri)
	}
}
