package persistence

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type URI string

func GetSnapshotInterval(parsedURI *url.URL) (snapshotInterval int, err error) {
	snapshotInterval = 3
	if snapshotIntervalStr := parsedURI.Query().Get("snapshotInterval"); snapshotIntervalStr != "" {
		if snapshotInterval, err = strconv.Atoi(snapshotIntervalStr); err != nil {
			return 0, fmt.Errorf("invalid snapshotInterval value: %s. %w", snapshotIntervalStr, err)
		}
	}
	return
}

// GetDBName returns the database name from the given URI.
// If no database name is specified, "memory" is returned.
// Example: "db:bbolt:./my-db?snapshotInterval=3" returns "bbolt".
func GetDBName(parsedURI *url.URL) (db string, err error) {
	db = "memory"
	if dbStr := parsedURI.Query().Get("db"); dbStr != "" {
		db = dbStr
	}
	return
}

func GetPath(parsedURI *url.URL) (path string, err error) {
	parts := strings.Split(parsedURI.Opaque, ":")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid persistence URI: %s", parsedURI.String())
	}

	return url.PathUnescape(parts[1])
}
