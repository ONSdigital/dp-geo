# dp-geo

Library for handling geographical areas, e.g. calculating bounding box from long/lat point

### geo package

Create a bounding box (circle) using the CircleToPolygon method on GeoPoint object. The Geo point object represents a longitude and latitude coordinate on Earth.

Instantiate a new geo object:

1. Use default object or create new object:

```Go
import "github.com/ONSdigital/dp-geo/geo"

...

    geo := geo.DefaultConfig
    // or
    geo = geo.New(numberOfworkers, numberOfMaximumSegments)
    // numberOfWorkers defines the number of maximum concurrent go routines used by methods 
    // on the receiver, "geo" and the number of maximum segments defines the upper boundary for
    // the number of coordinates generated to define a circle from a geopoint
```


Instantiate a new geopoint object:

```Go
import "github.com/ONSdigital/dp-geo/geo"

...

    geoPoint, err := geo.NewGeoPoint(lat, lon)
    if err != nil {
        // handle error
    }

...
```

Create Geo structure object based on new geo point above

```Go
    ...

    // radius is in metres and cannot be greater than the 
    // radius of the Earth, 6378137 metres defined by wgs84
    radius := 50
    segments = 20 // number of coordinates to create a bounding circle

    boundingArea, err := geo.CircleToPolygon(geoPoint, radius, segments)
    if err != nil {
        // handle error
    }

    // This will output an object defining the geometric type, in this case a polygon and an array of longitude and latitude coordinates with the first and last items of the array to be equal to close the bounding area.
    fmt.Printf("Bounding area: %v", boundingArea)
    ...
```

The bounding area/geo structure object created abides by [WKT format](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry) for geometric shapes.

Below is an example of the entire code structure

```Go
import "github.com/ONSdigital/dp-geo/geo"

...
const maxSegments = 200

type appConfig struct {
    ...
    geo *geo.Config
    ...
}

var (
    numberOfWorkers = 20
    numberOfMaximumSegments = 30
)

func main() {
    
    // Generate geo configuration
    appConfig.geo = geo.New(numberOfworkers, numberOfMaximumSegments)

    ...
}


func (c *appConfig) handleLatLon(lat, lon float64, radiusInMetres, segments int) error {
    ...
    geoPoint, err := appConfig.geo.NewGeoPoint(lat, lon)
    if err != nil {
        // handle error
    }

    boundingArea, err := geoPoint.CircleToPolygon(radiusInMetres, segments)
    if err != nil {
        // handle error
    }
    ...
}
```

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

