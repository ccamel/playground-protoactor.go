package provider

import (
	"fmt"
	"net/url"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	. "github.com/smartystreets/goconvey/convey" //nolint:revive,staticcheck

	providerv1 "github.com/ccamel/playground-protoactor.go/internal/persistence/provider/v1"
	"github.com/ccamel/playground-protoactor.go/internal/persistence/registry"
	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
)

func DoTest(uri string, factory registry.StoreFactory) {
	Convey("Given a store", func() {
		Convey("When creating a new store", func() {
			parsedURI, err := url.Parse(uri)
			So(err, ShouldBeNil)
			p, err := factory(nil, parsedURI)

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
			})

			nbEvents := uint8(15)
			const actorName = "test-actor"
			Convey(fmt.Sprintf("And when inserting a %d events", nbEvents), func() {
				records := make([]*persistencev1.EventRecord, nbEvents)

				for i := uint8(0); i < nbEvents; i++ {
					var event proto.Message
					if i%2 == 0 {
						event = &providerv1.SomethingHappened{Message: fmt.Sprintf("This is message %d", i)}
					} else {
						event = &providerv1.SomethingElseHappened{Value: uint64(i)}
					}

					payload, err := anypb.New(event)
					So(err, ShouldBeNil)

					id := fmt.Sprintf("%d", i)
					records[i] = &persistencev1.EventRecord{
						Id:        id,
						Type:      payload.TypeUrl,
						StreamId:  actorName,
						Version:   uint64(i),
						Timestamp: timestamppb.Now(),
						Payload:   payload,
					}
				}

				for _, record := range records {
					p.PersistEvent(actorName, record)
				}

				for version := uint8(0); version < nbEvents; version++ {
					Convey(fmt.Sprintf("Then all events are retrieved back from version %d", version), func() {
						count := version
						p.GetEvents(actorName, int(version), 0, func(record *persistencev1.EventRecord) {
							message, err := record.Payload.UnmarshalNew()
							So(err, ShouldBeNil)

							if count%2 == 0 {
								So(message, ShouldHaveSameTypeAs, &providerv1.SomethingHappened{})
							} else {
								So(message, ShouldHaveSameTypeAs, &providerv1.SomethingElseHappened{})
							}

							count++
						})

						So(count, ShouldEqual, nbEvents)
					})
				}
			})
		})
	})
}
