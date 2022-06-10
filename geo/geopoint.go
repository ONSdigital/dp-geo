package geo

import (
	"math"
	"sync"
)

// Coordinate describes the position of a point
type Coordinate struct {
	Lat float64
	Lon float64
}

// NewGeoPoint creates a coordinate object
func CreateGeoPoint(lat, lon float64) (*Coordinate, error) {
	geoPoint := &Coordinate{Lat: lat, Lon: lon}

	if err := geoPoint.validate(); err != nil {
		return nil, err
	}

	return geoPoint, nil
}

// CircleToPolygon generates a bounding box (circle) from a single geopoint using radius (in metres) and the number of segments
func (geo *Config) CircleToPolygon(geoPoint Coordinate, radius float64, segments int) (*GeoStructure, error) {

	// validate input
	if err := geo.validateInput(&geoPoint, radius, segments); err != nil {
		return nil, err
	}

	shape := &GeoStructure{
		Type: polygon,
	}

	coordinates := make([][]float64, segments)

	var semaphoreChan = make(chan struct{}, toLowest(segments, geo.defaultConcurrencyLimit))

	var wg sync.WaitGroup

	for i := 0; i < segments; i++ {

		semaphoreChan <- struct{}{}

		wg.Add(1)

		go func(i int) {
			sector := (twoPi * float64(-i)) / float64(segments)
			coordinate := generateCoordinate(geoPoint, radius, sector)
			coordinates[i] = coordinate

			<-semaphoreChan // read to release a slot
			wg.Done()
		}(i)
	}

	wg.Wait()

	// Push first coordinate to be last coordinate to complete polygon circle
	coordinates = append(coordinates, coordinates[0])

	shape.Coordinates = coordinates

	return shape, nil
}

func generateCoordinate(geoPoint Coordinate, distance float64, sector float64) []float64 {
	lat1 := toRadians(geoPoint.Lat)
	lon1 := toRadians(geoPoint.Lon)

	// distance divided by radius of the earth
	dByR := distance / radiusOfEarth

	lat := math.Asin(
		math.Sin(lat1)*math.Cos(dByR) + math.Cos(lat1)*math.Sin(dByR)*math.Cos(sector),
	)

	lon := lon1 + math.Atan2(
		math.Sin(sector)*math.Sin(dByR)*math.Cos(lat1),
		math.Cos(dByR)-math.Sin(lat1)*math.Sin(lat),
	)

	return []float64{toDegrees(lon), toDegrees(lat)}
}

func (geo *Config) validateInput(geoPoint *Coordinate, radius float64, segments int) error {

	if err := validateRadius(radius); err != nil {
		return err
	}

	if err := validateSegments(geo.defaultMaxSegments, segments); err != nil {
		return err
	}

	if err := geoPoint.validate(); err != nil {
		return err
	}

	return nil
}

func (geoPoint *Coordinate) validate() error {
	if geoPoint.Lon > 180 || geoPoint.Lon < -180 {
		return ErrInvalidLongitudinalPoint
	}

	if geoPoint.Lat > 90 || geoPoint.Lat < -90 {
		return ErrInvalidLatitudinalPoint
	}

	return nil
}
