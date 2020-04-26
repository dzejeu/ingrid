package webapp

import (
	"errors"
	"fmt"
	"ingrid/pkg/datarepr"
	"ingrid/pkg/externalapi"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func parseCoordinateString(cordsString string) (datarepr.GeoPoint, error) {
	cords := strings.SplitN(cordsString, ",", 2)
	long, err1 := strconv.ParseFloat(cords[0], 64)
	lat, err2 := strconv.ParseFloat(cords[1], 64)
	if err1 != nil || err2 != nil {
		return datarepr.GeoPoint{}, errors.New("unable to parse coordinates")
	} else {
		return datarepr.GeoPoint{Lat: lat, Long: long}, nil
	}
}

func extractRequestParams(request *http.Request) (datarepr.GeoPoint, []datarepr.GeoPoint, error){
	source := datarepr.GeoPoint{}
	var destinations []datarepr.GeoPoint
	for k, v := range request.URL.Query() {
		for _, s := range v {
			point, err := parseCoordinateString(s)
			if err != nil {
				return datarepr.GeoPoint{}, nil, err
			}

			if k == "src" {
				source = point
			} else if k == "dst" {
				destinations = append(destinations, point)
			}
		}
	}
	return source, destinations, nil
}

func HandleRequest(writer http.ResponseWriter, request *http.Request) {
	source, destinations, err := extractRequestParams(request)

	if err != nil {
		fmt.Fprintf(writer, fmt.Sprintf("error occured during parsing request params: %s", err.Error()))
	}

	var routes []datarepr.Route
	for _, dst := range destinations {
		route, err := externalapi.FetchRoute(source, dst)

		if err != nil {
			fmt.Fprintf(writer, fmt.Sprintf("error occured during data fetching: %s", err.Error()))
		} else {
			routes = append(routes, route)
		}
	}

	sort.Slice(routes, func(i int, j int) bool {
		route1 := routes[i]
		route2 := routes[j]
		if route1.Duration == route2.Duration {
			return route1.Distance < route2.Distance
		} else {
			return route1.Duration < route2.Duration
		}
	})

	writer.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(writer, datarepr.RoutingPlan{Source: source, Routes: routes}.JsonRepr())
}