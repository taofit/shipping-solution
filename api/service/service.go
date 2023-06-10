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

type Response struct {
	Code   string  `json:"code"`
	Routes []Route `json:"routes"`
}

func ValidateParameters(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query().Get("src")
	_, _, err := coordsparser.Parse(src)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dst := r.URL.Query()["dst"]
	if dst == nil {
		http.Error(w, "dst is empty", http.StatusBadRequest)
		return
	}
	for _, v := range dst {
		_, _, err := coordsparser.Parse(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func GetList(w http.ResponseWriter, src string, dst []string) (RouteResponse, error) {
	destinations, err := getDestinations(src, dst)
	if err != nil {
		return RouteResponse{}, err
	}
	destinationRoutes := sortDestinations(destinations)

	return RouteResponse{Source: src, Routes: destinationRoutes}, nil
}

func getDestinations(src string, dst []string) ([]DestinationRoute, error) {
	var response Response
	var destinationRoutes []DestinationRoute
	var destination DestinationRoute
	fmt.Println(dst)
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
		json.Unmarshal(body, &response)
		if strings.ToLower(response.Code) != "ok" {
			return nil, errors.New(response.Code)
		}
		destination = DestinationRoute{
			Destination: coor,
			Duration:    response.Routes[0].Duration,
			Distance:    response.Routes[0].Distance,
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
