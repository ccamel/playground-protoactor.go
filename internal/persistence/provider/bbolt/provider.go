package bbolt

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	bolt "go.etcd.io/bbolt"
)

type Provider struct {
	providerState persistence.ProviderState
}

func (p *Provider) GetState() persistence.ProviderState {
	return p.providerState
}

func NewProvider(system *actor.ActorSystem, path string, snapshotInterval int) (persistence.Provider, error) {
	db, err := bolt.Open(path, 0o666, nil)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("db", "bbolt").
		Str("path", db.Path()).
		Str("snapshotInterval", fmt.Sprintf("%d", snapshotInterval)).
		Msg("persistence provider started")

	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("events")); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("snapshots")); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Provider{
		providerState: &ProviderState{
			system:           system,
			snapshotInterval: snapshotInterval,
			db:               db,
			entropy:          ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0), //nolint:gosec
			subscribers:      &sync.Map{},
		},
	}, nil
}
