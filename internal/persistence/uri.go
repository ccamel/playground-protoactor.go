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
// Example: "db:bbolt:./my-db?snapshotInterval=3" returns "bbolt".
func GetDBName(parsedURI *url.URL) (string, error) {
	parts := strings.Split(parsedURI.Opaque, ":")

	db := parts[0]
	if db == "" {
		return "", fmt.Errorf("no database name found in URI: %s", parsedURI)
	}

	return db, nil
}

func GetPath(parsedURI *url.URL) (path string, err error) {
	parts := strings.Split(parsedURI.Opaque, ":")

	if len(parts) < 2 {
		return url.PathUnescape(parts[0])
	}

	return url.PathUnescape(parts[1])
}
