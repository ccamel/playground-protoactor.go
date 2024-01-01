package bbolt

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/rs/zerolog/log"
	bolt "go.etcd.io/bbolt"

	provider "github.com/ccamel/playground-protoactor.go/internal/persistence"
)

type Provider struct {
	providerState persistence.ProviderState
}

func (p *Provider) GetState() persistence.ProviderState {
	return p.providerState
}

func NewProvider(system *actor.ActorSystem, uri *url.URL) (persistence.Provider, error) {
	path, err := provider.GetPath(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid persistence URI: %s. %w", uri, err)
	}
	snapshotInterval, err := provider.GetSnapshotInterval(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid persistence URI: %s. %w", uri, err)
	}

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
			subscribers:      &sync.Map{},
		},
	}, nil
}

func init() {
	provider.RegisterFactory("db", NewProvider)
}
