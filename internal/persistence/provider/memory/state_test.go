package memory

import (
	"fmt"
	"net/url"
	"testing"

	"google.golang.org/protobuf/proto"

	. "github.com/smartystreets/goconvey/convey"

	memoryv1 "github.com/ccamel/playground-protoactor.go/internal/persistence/provider/memory/v1"
)

func TestNewProvider(t *testing.T) {
	Convey("Given a NewProvider", t, func() {
		uri, err := url.Parse(fmt.Sprintf("db:memory?snapshotInterval=%d", 5))
		So(err, ShouldBeNil)
		p, err := NewProvider(nil, uri)

		So(err, ShouldBeNil)
		So(p, ShouldNotBeNil)
		So(p.GetState(), ShouldNotBeNil)

		nbEvents := 15
		Convey(fmt.Sprintf("When inserting a %d events", nbEvents), func() {
			events := make([]proto.Message, nbEvents)

			for i := 0; i < nbEvents; i++ {
				if i%2 == 0 {
					events[i] = &memoryv1.SomethingHappened{Message: fmt.Sprintf("This is message %d", i)}
				} else {
					events[i] = &memoryv1.SomethingElseHappened{Value: uint64(i)}
				}
			}

			for i, event := range events {
				p.GetState().PersistEvent("test", i, event)
			}

			for version := 0; version < nbEvents; version++ {
				Convey(fmt.Sprintf("Then all events are retrieved back from version %d", version), func() {
					count := version
					p.GetState().GetEvents("test", version, 0, func(event interface{}) {
						if count%2 == 0 {
							So(event, ShouldHaveSameTypeAs, &memoryv1.SomethingHappened{})
							So(event.(*memoryv1.SomethingHappened).Message, ShouldEqual, events[count].(*memoryv1.SomethingHappened).Message)
						} else {
							So(event, ShouldHaveSameTypeAs, &memoryv1.SomethingElseHappened{})
							So(event.(*memoryv1.SomethingElseHappened).Value, ShouldEqual, events[count].(*memoryv1.SomethingElseHappened).Value)
						}

						count++
					})

					So(count, ShouldEqual, nbEvents)
				})
			}
		})
	})
}
