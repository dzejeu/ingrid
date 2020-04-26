package datarepr

import (
	"fmt"
	"strings"
)

type GeoPoint struct {
	Long float64
	Lat  float64
}

func (point GeoPoint) JsonRepr() string {
	return fmt.Sprintf("%f,%f", point.Long, point.Lat)
}

type Route struct {
	Destination GeoPoint
	Duration    float64
	Distance    float64
}

func (route Route) JsonRepr() string {
	return fmt.Sprintf(
		`{"destination": "%s", "duration": %f, "distance": %f}`,
		route.Destination.JsonRepr(),
		route.Duration,
		route.Distance)
}

type RoutingPlan struct {
	Source GeoPoint
	Routes []Route
}

func (routes RoutingPlan) JsonRepr() string {
	var routesAsJson []string
	for _, r := range routes.Routes {
		routesAsJson = append(routesAsJson, r.JsonRepr())
	}
	routesString := strings.Join(routesAsJson, ",")
	return fmt.Sprintf(`"source": "%s", "routes": [%s]`, routes.Source.JsonRepr(), routesString)
}