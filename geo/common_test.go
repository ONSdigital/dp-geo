package geo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestToLowest(t *testing.T) {

	Convey("Given 2 variables of type integer", t, func() {

		Convey("When passing the larger integer into toLowest first", func() {
			lowestValue := toLowest(30, 20)

			Convey("Then the return of function returns the lower integer", func() {
				So(lowestValue, ShouldEqual, 20)
			})
		})

		Convey("When passing the smaller integer into toLowest first", func() {
			lowestValue := toLowest(10, 40)

			Convey("Then the return of function returns the lower integer", func() {
				So(lowestValue, ShouldEqual, 10)
			})
		})
	})
}
