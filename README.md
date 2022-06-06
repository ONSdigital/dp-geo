# dp-geo

Library for handling geographical areas, e.g. calculating bounding box from long/lat point

### geo package

Create a bounding box (circle) using the CircleToPolygon method on GeoPoint object. The Geo point object represents a longitude and latitude coordinate on Earth.

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

    boundingArea, err := geoPoint.CircleToPolygon(radius, segments)
    if err != nil {
        // handle error
    }

    // This will output an object defining the geometric type, in this case a polygon and an array of longitude and latitude coordinates with the first and last items of the array to be equal to close the bounding area.
    fmt.Printf("Bounding area: %v", boundingArea)
    ...
```

The bounding area/geo structure object created abides by [WKT format](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry) for geometric shapes.

The number of coordinates generated for the bounding area is limited to the `defaultMaxSegments`. This can be changed to create a more precise circle, this is at a cost on how fast the polygon shape can be generated. To change the default maximum number of segments, call the `SetMaximumSegments` func

```Go
import "github.com/ONSdigital/dp-geo/geo"

...
const maxSegments = 200

func main() {
    
    // This will be used to provide the upper value of segments to be used in generating a polygon (circle)
    geo.SetMaximumSegments(MaxSegments)

    ...
}

func handleLatLon(lat, lon float64, radiusInMetres, segments int) error {
    ...
    geoPoint, err := geo.NewGeoPoint(lat, lon)
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

