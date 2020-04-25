package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Printable interface {
	asJson() string
}

type GeoPoint struct {
	Long float64
	Lat  float64
}

func (point GeoPoint) asJson() string {
	return fmt.Sprintf("%f,%f", point.Lat, point.Long)
}

type Route struct {
	Destination GeoPoint
	Duration    float64
	Distance    float64
}

func (route Route) asJson() string {
	return fmt.Sprintf(
		`{"destination": "%s", "duration": %f, "distance": %f}`,
		route.Destination.asJson(),
		route.Duration,
		route.Distance)
}

type SortedRoutes struct {
	Source GeoPoint `json:"source"`
	Routes []Route
}

func (routes SortedRoutes) asJson() string {
	var routesAsJson []string
	for _, r := range routes.Routes {
		routesAsJson = append(routesAsJson, r.asJson())
	}
	routesString := strings.Join(routesAsJson, ",")
	return fmt.Sprintf(`"source": "%s", "routes": [%s]`, routes.Source.asJson(), routesString)
}

type ApiResponse struct {
	Routes []map[string]interface{} `json:"routes"`
	Code   string                   `json:"code"`
}

func prepareApiUrl(src GeoPoint, dest GeoPoint) string {
	apiUri := "http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=false"
	return fmt.Sprintf(apiUri, src.Long, src.Lat, dest.Long, dest.Lat)
}

func fetchRoute(src GeoPoint, dest GeoPoint) Route {
	resp, err := http.Get(prepareApiUrl(src, dest))
	//TODO: extract proper fields

	if err != nil {
		// do something
		return Route{}
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var apiResp ApiResponse
		err := json.Unmarshal(body, &apiResp)
		fmt.Println(err) //TODO: handle this as well
		fmt.Println(apiResp.Code) //TODO: handle response code
		duration := apiResp.Routes[0]["duration"].(float64)
		distance := apiResp.Routes[0]["distance"].(float64)
		return Route{dest, duration, distance}
	}
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	source := GeoPoint{0, 0}
	var destinations []GeoPoint
	for k, v := range request.URL.Query() {
		for _, s := range v {
			cords := strings.SplitN(s, ",", 2)
			lat, _ := strconv.ParseFloat(cords[0], 64)
			long, _ := strconv.ParseFloat(cords[1], 64)
			if k == "src" {
				source.Lat = lat
				source.Long = long
			} else if k == "dst" {
				destinations = append(destinations, GeoPoint{lat, long})
			}
		}
	}
	var routes []Route
	for _, dst := range destinations {
		//TODO: gather all routes and sort
		routes = append(routes, fetchRoute(source, dst))
		//fmt.Fprintf(writer, "XD")
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
	fmt.Fprintf(writer, SortedRoutes{Source: source, Routes: routes}.asJson())
}

func main() {
	http.HandleFunc("/", handleRequest)

	http.ListenAndServe(":8080", nil)
}
