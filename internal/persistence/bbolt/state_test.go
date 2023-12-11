// Copyright © 2020 Chris Camel <camel.christophe@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package bbolt

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"google.golang.org/protobuf/proto"

	. "github.com/smartystreets/goconvey/convey"

	bboltv1 "github.com/ccamel/playground-protoactor.go/internal/persistence/bbolt/v1"
)

func TestNewProvider(t *testing.T) {
	Convey("Under temporary directory", t, func(c C) {
		dir, err := os.MkdirTemp("", "test-db-")
		So(err, ShouldBeNil)

		Convey("Given a NewProvider", func() {
			p, err := NewProvider(nil, path.Join(dir, "event-store.bolt.db"), 5)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.GetState(), ShouldNotBeNil)

			nbEvents := 15
			Convey(fmt.Sprintf("When inserting a %d events", nbEvents), func() {
				events := make([]proto.Message, nbEvents)

				for i := 0; i < nbEvents; i++ {
					if i%2 == 0 {
						events[i] = &bboltv1.SomethingHappened{Message: fmt.Sprintf("This is message %d", i)}
					} else {
						events[i] = &bboltv1.SomethingElseHappened{Value: uint64(i)}
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
								So(event, ShouldHaveSameTypeAs, &bboltv1.SomethingHappened{})
								So(event.(*bboltv1.SomethingHappened).Message, ShouldEqual, events[count].(*bboltv1.SomethingHappened).Message)
							} else {
								So(event, ShouldHaveSameTypeAs, &bboltv1.SomethingElseHappened{})
								So(event.(*bboltv1.SomethingElseHappened).Value, ShouldEqual, events[count].(*bboltv1.SomethingElseHappened).Value)
							}

							count++
						})

						So(count, ShouldEqual, nbEvents)
					})
				}
			})

			Reset(func() {
				err := json.NewEncoder(os.Stderr).Encode(p.GetState().(*ProviderState).db.Stats())
				So(err, ShouldBeNil)

				err = p.GetState().(*ProviderState).Close()
				So(err, ShouldBeNil)

				err = os.RemoveAll(dir)
				So(err, ShouldBeNil)
			})
		})
	})
}
