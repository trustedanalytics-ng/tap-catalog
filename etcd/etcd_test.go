// +build dev_test
package etcd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEtcdKVStore(t *testing.T) {
	client := EtcdConnector{}
	Convey("Test EtcdKVStore", t, func() {
		Convey("Should update properly", func() {
			err := client.Set("test", "testValue", "", 0)
			So(err, ShouldBeNil)
		})

		Convey("Should get KV properly", func() {
			value, err := client.GetKeyValue("test")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "testValue")
		})

		Convey("Should get KV into struct properly", func() {
			result := ""
			err := client.GetKeyIntoStruct("test", &result)
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "testValue")
		})

		Convey("Should delete KV properly", func() {
			err := client.Delete("test", 0)
			So(err, ShouldBeNil)
		})
	})
}
