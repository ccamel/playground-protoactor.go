package registry

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asynkron/protoactor-go/actor"

	persistence2 "github.com/ccamel/playground-protoactor.go/internal/persistence"
)

type (
	StoreFactory  func(system *actor.ActorSystem, uri *url.URL) (persistence2.Store, error)
	StoreRegistry map[string]StoreFactory
)

func (f StoreRegistry) Get(name string) (StoreFactory, error) {
	factory, ok := f[name]
	if !ok {
		return nil, fmt.Errorf(
			"unsupported persistence scheme: %s. Supported schemes: %s",
			name,
			strings.Join(SupportedSchemes(), ", "))
	}
	return factory, nil
}

func (f StoreRegistry) GetFromURI(uri *url.URL) (StoreFactory, error) {
	db, err := persistence2.GetDBName(uri)
	if err != nil {
		return nil, err
	}
	return f.Get(db)
}

func SupportedSchemes() []string {
	schemes := make([]string, 0, len(factories))
	for scheme := range factories {
		schemes = append(schemes, scheme)
	}
	return schemes
}

// factories is the list of registered stores.
var factories StoreRegistry = make(map[string]StoreFactory)

// RegisterFactory registers a store factory given its name.
func RegisterFactory(name string, factory func(system *actor.ActorSystem, uri *url.URL) (persistence2.Store, error)) error {
	if _, exists := factories[name]; exists {
		return fmt.Errorf("store factory already registered: %s", name)
	}

	factories[name] = factory
	return nil
}
