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
package usr

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ccamel/playground-protoactor.go/internal/actor/booklend"
	"github.com/golang/protobuf/ptypes"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/code"
)

type Actor struct {
}

func (state *Actor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		pid, _ := context.SpawnNamed(booklend.NewBookCommandHandler(), "book_lend")

		bookID := "207b6be6-a7e4-4cc7-a692-b51e79de0460" // uuid.New().String()

		res, err := context.RequestFuture(pid, &booklend.RegisterBook{
			BookId: bookID,
			Title:  "The Lord of the Rings",
			Isbn:   "0-618-15396-9",
		}, 5*time.Second).Result()

		if err != nil {
			log.Error().Err(err).Msg("Failed to create book")

			return
		}

		if res.(*booklend.CommandStatus).Code != code.Code_OK {
			log.Warn().Interface("event", res).Msgf("error")

			return
		}

		log.Info().Interface("event", res).Msgf("ok")

		res, err = context.RequestFuture(pid, &booklend.LendBook{
			BookId:           bookID,
			Borrower:         "John Doe",
			Date:             ptypes.TimestampNow(),
			ExpectedDuration: ptypes.DurationProto(90 * 24 * time.Hour),
		}, 5*time.Second).Result()

		if err != nil {
			log.Error().Err(err).Msg("Failed to lend book")

			return
		}

		if res.(*booklend.CommandStatus).Code != code.Code_OK {
			log.Warn().Interface("event", res).Msgf("error")

			return
		}

		log.Info().Interface("event", res).Msgf("ok")

		res, err = context.RequestFuture(pid, &booklend.ReturnBook{
			BookId: bookID,
			Date:   ptypes.TimestampNow(),
		}, 5*time.Second).Result()

		if err != nil {
			log.Error().Err(err).Msg("Failed to return book")

			return
		}

		if res.(*booklend.CommandStatus).Code != code.Code_OK {
			log.Warn().Interface("event", res).Msgf("error")

			return
		}

		log.Info().Interface("event", res).Msgf("ok")
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	}
}

func New() actor.Producer {
	return func() actor.Actor {
		return &Actor{}
	}
}
