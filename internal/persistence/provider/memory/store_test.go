package memory

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ccamel/playground-protoactor.go/internal/persistence/provider"
)

func TestNewProvider(t *testing.T) {
	Convey("Given a NewStore", t, func() {
		uri := fmt.Sprintf("db:memory?snapshotInterval=%d", 5)

		provider.DoTest(uri, NewStore)
	})
}
