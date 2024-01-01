package persistence

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
)

type (
	ProviderFactory func(system *actor.ActorSystem, uri *url.URL) (persistence.Provider, error)
	FactoryRegistry map[string]ProviderFactory
)

func (f FactoryRegistry) Get(name string) (ProviderFactory, error) {
	factory, ok := f[name]
	if !ok {
		return nil, fmt.Errorf(
			"unsupported persistence scheme: %s. Supported schemes: %s",
			name,
			strings.Join(getSupportedSchemes(), ", "))
	}
	return factory, nil
}

func (f FactoryRegistry) GetFromURI(uri *url.URL) (ProviderFactory, error) {
	db, err := GetDBName(uri)
	if err != nil {
		return nil, err
	}
	return f.Get(db)
}

func getSupportedSchemes() []string {
	schemes := make([]string, 0, len(factories))
	for scheme := range factories {
		schemes = append(schemes, scheme)
	}
	return schemes
}

// factories is the list of registered persistence providers.
var factories FactoryRegistry = make(map[string]ProviderFactory)

func RegisterFactory(name string, factory func(system *actor.ActorSystem, uri *url.URL) (persistence.Provider, error)) {
	factories[name] = factory
}

func NewProvider(system *actor.ActorSystem, uri URI) (persistence.Provider, error) {
	if uri == "" {
		return nil, fmt.Errorf("persistence URI is required")
	}

	parsedURI, err := url.Parse(string(uri))
	if err != nil {
		return nil, err
	}

	if parsedURI.Scheme != "db" {
		return nil, fmt.Errorf("invalid persistence URI: %s", uri)
	}

	factory, err := factories.Get(parsedURI.Scheme)
	if err != nil {
		return nil, err
	}

	return factory(system, parsedURI)
}
