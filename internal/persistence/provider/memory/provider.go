package memory

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/rs/zerolog/log"

	provider "github.com/ccamel/playground-protoactor.go/internal/persistence"
	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
)

type Provider struct {
	providerState persistence.ProviderState
}

func (p *Provider) GetState() persistence.ProviderState {
	return p.providerState
}

func NewProvider(system *actor.ActorSystem, uri *url.URL) (persistence.Provider, error) {
	snapshotInterval, err := provider.GetSnapshotInterval(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid persistence URI: %s. %w", uri, err)
	}

	log.Info().
		Str("db", "memory").
		Str("snapshotInterval", fmt.Sprintf("%d", snapshotInterval)).
		Msg("persistence provider started")

	return &Provider{
		providerState: &ProviderState{
			system:           system,
			events:           make(map[string][]*persistencev1.EventRecord),
			snapshots:        make(map[string]*persistencev1.SnapshotRecord),
			snapshotInterval: snapshotInterval,
			subscribers:      &sync.Map{},
		},
	}, nil
}

func init() {
	provider.RegisterFactory("db", NewProvider)
}
