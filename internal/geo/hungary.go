package geo

import (
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

var polygons []orb.Polygon

func LoadHungary(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return err
	}

	for _, f := range fc.Features {

		switch g := f.Geometry.(type) {

		case orb.Polygon:
			polygons = append(polygons, g)

		case orb.MultiPolygon:
			for _, p := range g {
				polygons = append(polygons, p)
			}
		}
	}

	return nil
}

func Contains(lat, lng float64) bool {

	p := orb.Point{lng, lat}

	for _, poly := range polygons {
		if planar.PolygonContains(poly, p) {
			return true
		}
	}

	return false
}
