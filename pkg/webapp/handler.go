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
	"time"
)

type response struct {
	route datarepr.Route
	err error
}

// Parser coordinate string which expected format is as "float,float" - ex. "13.2746,52.8876"
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

// Extracts params from http request url. Expected keys are "src" for source and "dst" for destination
// "dst" can appear multiple times
// Request example: http://app/?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219
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

// Fetches routes' information from third party api in asynchronous manner and responds with json formatted output
func HandleRequest(writer http.ResponseWriter, request *http.Request) {
	source, destinations, err := extractRequestParams(request)

	if err != nil {
		fmt.Fprintf(writer, fmt.Sprintf("error occured during parsing request params: %s", err.Error()))
	}

	var routes []datarepr.Route
	ch := make(chan *response)
	for _, dst := range destinations {
		go func(source datarepr.GeoPoint, dst datarepr.GeoPoint) {
			route, err := externalapi.FetchRoute(source, dst)
			ch <- &response{route: route, err: err}
		}(source, dst)
	}

loop:
	for {
		select {
		case r := <-ch:
			if (*r).err != nil {
				fmt.Fprintf(writer, fmt.Sprintf("error occured during data fetching: %s", (*r).err.Error()))
			}
			routes = append(routes, (*r).route)
			if len(routes) == len(destinations) {
				break loop
			}
		case <-time.After(50 * time.Millisecond):
			continue loop
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