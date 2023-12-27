package bbolt

import "github.com/asynkron/protoactor-go/persistence"

type Provider struct {
	providerState persistence.ProviderState
}

func (p *Provider) GetState() persistence.ProviderState {
	return p.providerState
}
