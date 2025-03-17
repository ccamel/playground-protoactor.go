package bbolt

import (
	"fmt"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ccamel/playground-protoactor.go/internal/persistence/provider"
)

func TestStore(t *testing.T) {
	Convey("Under temporary directory", t, func(_ C) {
		dir := t.TempDir()

		Convey("Given a NewStore", func() {
			uri := fmt.Sprintf("db:bbolt:%s?snapshotInterval=%d", path.Join(dir, "event-store.bolt.db"), 5)

			provider.DoTest(uri, NewStore)
		})
	})
}
