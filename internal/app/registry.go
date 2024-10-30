package app

import (
	"errors"
	"iter"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
)

var (
	errAlreadyRegistered  = errors.New("factory already registered")
	errNameCannotBeEmpty  = errors.New("name cannot be empty")
	errFactoryCannotBeNil = errors.New("factory cannot be nil")
)

var (
	registry Registry
	rwMu     sync.RWMutex
)

type (
	// Registry maps application names to their factory functions.
	Registry map[string]FactoryFunc
	// FactoryFunc creates an actor producer for a given configuration.
	// The configuration is an arbitrary value that is passed to the factory function and that is used to
	// create the application actor.
	FactoryFunc func(config any) actor.Producer
)

func init() {
	registry = make(Registry)
}

// Register registers the given factory with the given name.
// It returns an error if a factory with the same name is already registered, if the name is empty or
// if the factory is nil.
func Register(name string, factory FactoryFunc) error {
	if name == "" {
		return errNameCannotBeEmpty
	}
	if factory == nil {
		return errFactoryCannotBeNil
	}

	rwMu.Lock()
	defer rwMu.Unlock()

	if _, exists := registry[name]; exists {
		return errAlreadyRegistered
	}
	registry[name] = factory

	return nil
}

// RegisterMust registers the given factory with the given name.
// It panics if a factory with the same name is already registered, if the name is empty or
// if the factory is nil.
func RegisterMust(name string, factory FactoryFunc) {
	if err := Register(name, factory); err != nil {
		panic(err)
	}
}

// Seq returns a function that iterates over the registry with a callback.
func Seq() iter.Seq2[string, FactoryFunc] {
	return func(yield func(string, FactoryFunc) bool) {
		rwMu.RLock()
		defer rwMu.RUnlock()

		for k, v := range registry {
			if !yield(k, v) {
				return
			}
		}
	}
}
