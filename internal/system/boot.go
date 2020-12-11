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
package system

import (
	SYS "syscall"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/actor/middleware/propagator"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/ccamel/playground-protoactor.go/internal/middleware"
	"github.com/ccamel/playground-protoactor.go/internal/persistence/bbolt"
	"github.com/ccamel/playground-protoactor.go/internal/system/core"
	"github.com/rs/zerolog/log"
	DEATH "github.com/vrecan/death"
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

func Boot() (*System, error) {
	log.Info().
		Str("actor", "/").
		Msg("booting the system...")

	system := actor.NewActorSystem()

	log.Info().
		Str("actor", "/").
		Msg("start remote server...")

	log.Info().
		Str("actor", "/").
		Msg("remote server started")

	provider, err := bbolt.NewProvider(system, "my-db", 3)
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

	props := actor.PropsFromProducer(core.New()).WithSupervisor(actor.RestartingSupervisorStrategy())

	pid, err := rootContext.SpawnNamed(props, "init")
	if err != nil {
		return nil, err
	}

	return &System{
		rootContext: rootContext,
		initPid:     pid,
	}, nil
}
