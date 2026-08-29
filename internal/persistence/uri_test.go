package persistence

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDBName(t *testing.T) {
	Convey("Given an opaque URI without a database delimiter", t, func() {
		dbName, err := GetDBName(&url.URL{Opaque: "path"})

		Convey("It returns the missing database name error", func() {
			So(dbName, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no database name found in URI: path")
		})
	})

	Convey("Given an opaque URI with a database name", t, func() {
		dbName, err := GetDBName(&url.URL{Opaque: "bbolt:./my-db"})

		Convey("It returns the database name", func() {
			So(dbName, ShouldEqual, "bbolt")
			So(err, ShouldBeNil)
		})
	})
}
