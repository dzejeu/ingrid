package externalapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"ingrid/pkg/datarepr"
	"io/ioutil"
	"net/http"
)

type apiResponse struct {
	Routes []map[string]interface{} `json:"routes"`
	Code   string                   `json:"code"`
}

func prepareApiUrl(src datarepr.GeoPoint, dest datarepr.GeoPoint) string {
	apiUri := "http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=false"
	return fmt.Sprintf(apiUri, src.Long, src.Lat, dest.Long, dest.Lat)
}

func FetchRoute(src datarepr.GeoPoint, dest datarepr.GeoPoint) (datarepr.Route, error) {
	resp, err := http.Get(prepareApiUrl(src, dest))

	if err != nil {
		return datarepr.Route{}, err
	} else {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		var apiResp apiResponse
		jsonErr := json.Unmarshal(body, &apiResp)
		if jsonErr != nil {
			return datarepr.Route{}, jsonErr
		}

		if apiResp.Code != "Ok" {
			return datarepr.Route{}, errors.New(fmt.Sprintf("Unexpected response code: %s", apiResp.Code))
		}
		duration := apiResp.Routes[0]["duration"].(float64)
		distance := apiResp.Routes[0]["distance"].(float64)
		return datarepr.Route{Destination: dest, Duration: duration, Distance: distance}, nil
	}
}