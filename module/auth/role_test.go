package auth

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRole(t *testing.T) {

	r := Role{
		Id:   "test-1",
		Name: "Test",
		Desc: "a test role",
		Allow: []string{
			"*.abc.com/allow/*",
			"*.allow.com/*",
		},
		Deny: []string{
			"*.abc.com/deny/*",
		},
	}
	r.Init()

	Convey("Access to different urls", t, func() {

		Convey("Deny to access", func() {
			So(r.HasPermission("http://hi.abc.com/deny/you.html"), ShouldBeFalse)
			So(r.HasPermission("http://unset.com/you.html"), ShouldBeFalse)
		})
		Convey("Allow to access", func() {
			So(r.HasPermission("http://hi.abc.com/allow/you.html"), ShouldBeTrue)
			So(r.HasPermission("http://www.allow.com/you.html"), ShouldBeTrue)
		})
	})
}
