package persistence

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"

	"github.com/ccamel/playground-protoactor.go/internal/persistence/provider/bbolt"
	"github.com/ccamel/playground-protoactor.go/internal/persistence/provider/memory"
)

type URI string

func NewProvider(system *actor.ActorSystem, uri URI) (persistence.Provider, error) {
	if uri == "" {
		return nil, fmt.Errorf("persistence URI is required")
	}

	parsedURI, err := url.Parse(string(uri))
	if err != nil {
		return nil, err
	}

	switch parsedURI.Scheme {
	case "db":
		parts := strings.Split(parsedURI.Opaque, ":")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid persistence URI: %s", uri)
		}
		head := parts[0]
		tail := parts[1]

		switch head {
		case "bbolt":
			path, err := url.PathUnescape(tail)
			if err != nil {
				return nil, fmt.Errorf("invalid persistence URI: %s. %w", uri, err)
			}

			snapshotInterval, err := getSnapshotInterval(parsedURI)
			if err != nil {
				return nil, fmt.Errorf("invalid persistence URI: %s. %w", uri, err)
			}

			return bbolt.NewProvider(system, path, snapshotInterval)
		case "memory":
			snapshotInterval, err := getSnapshotInterval(parsedURI)
			if err != nil {
				return nil, fmt.Errorf("invalid persistence URI: %s. %w", uri, err)
			}

			return memory.NewProvider(system, snapshotInterval)
		default:
			return nil, fmt.Errorf("unsupported database: %s", head)
		}
	default:
		return nil, fmt.Errorf("invalid persistence URI: %s", uri)
	}
}

func getSnapshotInterval(parsedURI *url.URL) (snapshotInterval int, err error) {
	snapshotInterval = 3
	if snapshotIntervalStr := parsedURI.Query().Get("snapshotInterval"); snapshotIntervalStr != "" {
		if snapshotInterval, err = strconv.Atoi(snapshotIntervalStr); err != nil {
			return 0, fmt.Errorf("invalid snapshotInterval value: %s. %w", snapshotIntervalStr, err)
		}
	}
	return
}
