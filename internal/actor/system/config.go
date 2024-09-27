package system

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ccamel/playground-protoactor.go/internal/persistence"
)

type Config struct {
	PersistenceURI persistence.URI
}

type Option func(*Config) error

// NewConfig builds a new config using the given options.
func NewConfig(opts ...Option) (*Config, error) {
	c := &Config{}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

var (
	ErrURIParseError = errors.New("URI parse error")
	ErrInvalidScheme = errors.New("invalid scheme")
)

// WithPersistenceURI sets the persistence URI.
// Example: "db:bbolt:./my-db".
func WithPersistenceURI(uri string) Option {
	return func(c *Config) error {
		u, err := url.Parse(uri)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrURIParseError, err)
		}
		if u.Scheme != "db" {
			return fmt.Errorf("%w: scheme must be 'db'", ErrInvalidScheme)
		}

		c.PersistenceURI = persistence.URI(uri)

		return nil
	}
}
