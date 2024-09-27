package memory

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/rs/zerolog/log"

	persistence2 "github.com/ccamel/playground-protoactor.go/internal/persistence"
	"github.com/ccamel/playground-protoactor.go/internal/persistence/registry"
	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
)

const DBName = "memory"

func NewStore(system *actor.ActorSystem, uri *url.URL) (persistence2.Store, error) {
	snapshotInterval, err := persistence2.GetSnapshotInterval(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid persistence URI: %s. %w", uri, err)
	}

	log.Info().
		Str("db", DBName).
		Str("snapshotInterval", fmt.Sprintf("%d", snapshotInterval)).
		Msg("persistence provider started")

	return &Store{
		system:           system,
		events:           make(map[string][]*persistencev1.EventRecord),
		snapshots:        make(map[string]*persistencev1.SnapshotRecord),
		snapshotInterval: snapshotInterval,
		subscribers:      &sync.Map{},
	}, nil
}

func init() {
	if err := registry.RegisterFactory(DBName, NewStore); err != nil {
		panic(err)
	}
}
