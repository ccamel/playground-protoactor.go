package entropy

import (
	"crypto/rand"
	"io"
	"sync"

	"github.com/oklog/ulid/v2"
)

var defaultEntropy = sync.OnceValue(func() io.Reader {
	rng := ulid.Monotonic(rand.Reader, 0)

	return &ulid.LockedMonotonicReader{
		MonotonicReader: ulid.Monotonic(rng, 0),
	}
})

// MakeULID returns an ULID with the current time in Unix milliseconds and
// monotonically increasing entropy for the same millisecond.
// It is safe for concurrent use.
func MakeULID() ulid.ULID {
	// NOTE: MustNew can't panic since defaultEntropy() never returns an error.
	return ulid.MustNew(ulid.Now(), defaultEntropy())
}
