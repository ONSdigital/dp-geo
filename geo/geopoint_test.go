package geo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateGeoPoint(t *testing.T) {
	Convey("Given a valid latitudinal and longitudinal point", t, func() {
		lat := float64(90)
		lon := float64(180)

		Convey("When calling NewGeoPoint", func() {
			expectedGeoPoint := &Coordinate{
				Lat: lat,
				Lon: lon,
			}

			geoPoint, err := CreateGeoPoint(lat, lon)

			Convey("Then geo point object is returned", func() {
				So(err, ShouldBeNil)
				So(geoPoint, ShouldResemble, expectedGeoPoint)
			})
		})
	})

	Convey("Given an invalid latitudinal or longitudinal point", t, func() {
		lat := -90.121
		lon := 25.34343

		Convey("When calling NewGeoPoint", func() {

			geoPoint, err := CreateGeoPoint(lat, lon)

			Convey("Then an error is returned", func() {
				So(err, ShouldEqual, ErrInvalidLatitudinalPoint)
				So(geoPoint, ShouldBeNil)
			})
		})
	})
}

func TestCoordinateValidation(t *testing.T) {
	Convey("Given a valid latitudinal and longitudinal point", t, func() {
		lat := float64(-90)
		lon := float64(-180)

		validCoordinate := &Coordinate{
			Lat: lat,
			Lon: lon,
		}

		Convey("When calling validate", func() {

			err := validCoordinate.validate()

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given an invalid latitudinal point", t, func() {
		Convey("And it is less than the lower bound value (-90)", func() {

			lat := -90.0000001
			lon := -25.34343

			coordinate := &Coordinate{
				Lat: lat,
				Lon: lon,
			}

			Convey("When calling validate", func() {

				err := coordinate.validate()

				Convey("Then no error is returned", func() {
					So(err, ShouldEqual, ErrInvalidLatitudinalPoint)
				})
			})
		})

		Convey("And it is greater than the upper bound value (90)", func() {

			lat := 90.0000001
			lon := -25.34343

			coordinate := &Coordinate{
				Lat: lat,
				Lon: lon,
			}

			Convey("When calling validate", func() {

				err := coordinate.validate()

				Convey("Then no error is returned", func() {
					So(err, ShouldEqual, ErrInvalidLatitudinalPoint)
				})
			})
		})
	})

	Convey("Given an invalid longitudinal point", t, func() {
		Convey("And it is less than the lower bound value (-180)", func() {

			lat := 42.3333
			lon := 180.0000001

			coordinate := &Coordinate{
				Lat: lat,
				Lon: lon,
			}

			Convey("When calling validate", func() {

				err := coordinate.validate()

				Convey("Then no error is returned", func() {
					So(err, ShouldEqual, ErrInvalidLongitudinalPoint)
				})
			})
		})

		Convey("And it is greater than the upper bound value (180)", func() {

			lat := 75.45
			lon := -180.0000001

			coordinate := &Coordinate{
				Lat: lat,
				Lon: lon,
			}

			Convey("When calling validate", func() {

				err := coordinate.validate()

				Convey("Then no error is returned", func() {
					So(err, ShouldEqual, ErrInvalidLongitudinalPoint)
				})
			})
		})
	})
}

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

func TestValidateInput(t *testing.T) {
	geo := DefaultConfig

	validCoordinate := &Coordinate{
		Lat: 23.4567,
		Lon: -34.765322,
	}

	Convey("Given coordinates, radius and segment values are valid", t, func() {

		Convey("When calling validateInput", func() {

			err := geo.validateInput(validCoordinate, 50, 10)

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given invalid inputs", t, func() {
		Convey("where the coordinates are invalid but radius and segments are not", func() {
			invalidCoordinate := &Coordinate{
				Lat: -90.1234,
				Lon: 23.435352,
			}

			Convey("When calling validateInput", func() {

				err := geo.validateInput(invalidCoordinate, 50, 10)

				Convey("Then an error is returned", func() {
					So(err, ShouldEqual, ErrInvalidLatitudinalPoint)
				})
			})
		})

		Convey("where the radius is invalid but coordinates and segments are not", func() {
			Convey("When calling validateInput", func() {

				err := geo.validateInput(validCoordinate, 987654321, 10)

				Convey("Then an error is returned", func() {
					So(err, ShouldEqual, ErrRadiusLargerThanEarth)
				})
			})
		})

		Convey("where the segments is invalid but coordinates and radius are not", func() {
			Convey("When calling validateInput", func() {

				err := geo.validateInput(validCoordinate, 100, 2)

				Convey("Then an error is returned", func() {
					So(err, ShouldEqual, ErrTooFewSegments)
				})
			})
		})
	})
}

func TestGenerateCoordinate(t *testing.T) {

	Convey("Given a valid coordinate, radius and sector of a circle", t, func() {
		coordinate := Coordinate{
			Lat: 23.4567,
			Lon: -34.765322,
		}

		sector := (twoPi * float64(1)) / float64(20)

		Convey("When calling generateCoordinate", func() {
			expectedCoordinate := []float64{-34.76517069885503, 23.457127174229374}
			newCoordinate := generateCoordinate(coordinate, 50, sector)

			Convey("Then a new coordinate is created relative to the original coordinate", func() {
				So(newCoordinate, ShouldResemble, expectedCoordinate)
			})
		})
	})
}

func TestCircleToPolygon(t *testing.T) {
	geo := DefaultConfig

	Convey("Given a valid geo point, radius and number of segments", t, func() {
		coordinate := Coordinate{
			Lat: 23.4567,
			Lon: -34.765322,
		}

		Convey("When calling generateCoordinate", func() {
			expectedShape := &GeoStructure{
				Type: polygon,
				Coordinates: [][]float64{{-34.765322000000005, 23.457149157642057},
					{-34.76560979174064, 23.457063375901644},
					{-34.76578765602604, 23.456838796653578},
					{-34.7657876550471, 23.456561201964476},
					{-34.76560979015668, 23.456336623570493},
					{-34.765322000000005, 23.456250842357942},
					{-34.76503420984332, 23.456336623570493},
					{-34.764856344952904, 23.456561201964476},
					{-34.76485634397396, 23.456838796653578},
					{-34.765034208259365, 23.457063375901644},
					{-34.765322000000005, 23.457149157642057}},
			}

			shape, err := geo.CircleToPolygon(coordinate, 50, 10)

			Convey("Then a new polygon shape is created encircling the geopoint", func() {
				So(err, ShouldBeNil)
				So(shape.Coordinates[0], ShouldResemble, shape.Coordinates[10])
				So(shape, ShouldResemble, expectedShape)
			})
		})
	})
}

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
