package bbolt

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/rs/zerolog/log"
	bolt "go.etcd.io/bbolt"

	persistence2 "github.com/ccamel/playground-protoactor.go/internal/persistence"
	"github.com/ccamel/playground-protoactor.go/internal/persistence/registry"
)

const DBName = "bbolt"

func NewStore(system *actor.ActorSystem, uri *url.URL) (persistence2.Store, error) {
	path, err := persistence2.GetPath(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid persistence URI: %s. %w", uri, err)
	}
	snapshotInterval, err := persistence2.GetSnapshotInterval(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid persistence URI: %s. %w", uri, err)
	}

	db, err := bolt.Open(path, 0o666, nil)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("db", DBName).
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

	return &Store{
		system:           system,
		snapshotInterval: snapshotInterval,
		db:               db,
		subscribers:      &sync.Map{},
	}, nil
}

func init() {
	if err := registry.RegisterFactory(DBName, NewStore); err != nil {
		panic(err)
	}
}
