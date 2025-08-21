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

//nolint:funlen,gocognit
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

				Convey("Then unbounded iteration from any start version streams all remaining events", func() {
					for start := uint8(0); start < nbEvents; start++ {
						Convey(fmt.Sprintf("From %d to end", start), func() {
							count := start
							p.GetEvents(actorName, int(start), 0, func(record *persistencev1.EventRecord) {
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

				Convey("And then bounded ranges work correctly", func() {
					cases := []struct {
						start, end int
						expected   []uint64
					}{
						{3, 7, []uint64{3, 4, 5, 6, 7}},
						{5, 5, []uint64{5}},
						{10, 12, []uint64{10, 11, 12}},
						{8, 6, nil},
						{0, 0, []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}},
						{-5, 7, []uint64{0, 1, 2, 3, 4, 5, 6, 7}},
						{3, -1, nil},
					}

					for _, tc := range cases {
						Convey(fmt.Sprintf("Range [%d,%d] yields %v", tc.start, tc.end, tc.expected), func() {
							var got []uint64
							p.GetEvents(actorName, tc.start, tc.end, func(record *persistencev1.EventRecord) {
								got = append(got, record.Version)
							})
							So(got, ShouldResemble, tc.expected)
						})
					}
				})
			})
		})
	})
}
