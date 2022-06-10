package geo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateSegments(t *testing.T) {
	maxSegments := 180
	Convey("Given segments are within range of allowed values", t, func() {

		Convey("When calling validateSegments at the lower bound value (3)", func() {

			err := validateSegments(maxSegments, 3)

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When calling validateSegments at the upper bound value (180)", func() {

			err := validateSegments(maxSegments, 180)

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given segments are outside range of allowed values", t, func() {

		Convey("When calling validateSegments below the lower bound value of 3", func() {

			err := validateSegments(maxSegments, 2)

			Convey("Then an error is returned", func() {
				So(err, ShouldEqual, ErrTooFewSegments)
			})
		})

		Convey("When calling validateSegments above the upper bound value of 180", func() {

			err := validateSegments(maxSegments, 181)

			Convey("Then an error is returned", func() {
				So(err, ShouldResemble, ErrTooManySegments(180))
			})
		})
	})
}
