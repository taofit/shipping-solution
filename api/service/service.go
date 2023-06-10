package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/flopp/go-coordsparser"
)

type RouteResponse struct {
	Source string             `json:"source"`
	Routes []DestinationRoute `json:"routes"`
}

type DestinationRoute struct {
	Destination string  `json:"destination"`
	Duration    float32 `json:"duration"`
	Distance    float32 `json:"distance"`
}

type Route struct {
	Duration float32 `json:"duration"`
	Distance float32 `json:"distance"`
}

type ThirdPartyResponse struct {
	Code   string  `json:"code"`
	Routes []Route `json:"routes"`
}

func ValidateParameters(src []string, dst []string) string {
	if src == nil {
		return "source location is empty"
	}
	if len(src) > 1 {
		return "only one source location can be provided"
	}
	_, _, err := coordsparser.Parse(src[0])
	if err != nil {
		return err.Error()
	}

	if dst == nil {
		return "destination location is empty"
	}

	for _, v := range dst {
		_, _, err := coordsparser.Parse(v)
		if err != nil {
			return err.Error()
		}
	}

	return ""
}

func GetList(src string, dst []string) (RouteResponse, error) {
	destinations, err := getDestinations(src, dst)
	if err != nil {
		return RouteResponse{}, err
	}
	destinationRoutes := sortDestinations(destinations)

	return RouteResponse{Source: src, Routes: destinationRoutes}, nil
}

func getDestinations(src string, dst []string) ([]DestinationRoute, error) {
	var thirdPartyResponse ThirdPartyResponse
	var destinationRoutes []DestinationRoute
	var destination DestinationRoute

	for _, coor := range dst {
		url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%s;%s?overview=false", src, coor)
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(body, &thirdPartyResponse)
		if strings.ToLower(thirdPartyResponse.Code) != "ok" {
			return nil, errors.New(thirdPartyResponse.Code)
		}
		destination = DestinationRoute{
			Destination: coor,
			Duration:    thirdPartyResponse.Routes[0].Duration,
			Distance:    thirdPartyResponse.Routes[0].Distance,
		}
		destinationRoutes = append(destinationRoutes, destination)
	}

	return destinationRoutes, nil
}

func sortDestinations(destinationRoutes []DestinationRoute) []DestinationRoute {
	sort.Slice(destinationRoutes, func(i, j int) bool {
		if destinationRoutes[i].Duration == destinationRoutes[j].Duration {
			return destinationRoutes[i].Destination < destinationRoutes[j].Destination
		}
		return destinationRoutes[i].Duration < destinationRoutes[j].Duration
	})

	return destinationRoutes
}
