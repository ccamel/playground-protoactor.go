package memory

import (
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/rs/zerolog/log"

	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
)

type Provider struct {
	providerState persistence.ProviderState
}

func (p *Provider) GetState() persistence.ProviderState {
	return p.providerState
}

func NewProvider(system *actor.ActorSystem, snapshotInterval int) (persistence.Provider, error) {
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
