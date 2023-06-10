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

type DestinationRoute struct {
	destination string
	duration    float32
	distance    float32
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

func GetList(w http.ResponseWriter, src string, dst []string) {
	destinations, err := getDestinations(src, dst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	destinationRoutes := sortDestinations(destinations)
	fmt.Println(destinationRoutes)
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
		// fmt.Println(body)
		json.Unmarshal(body, &response)
		if strings.ToLower(response.Code) != "ok" {
			return nil, errors.New(response.Code)
		}
		// fmt.Println(response, response.Code, response.Routes[0].Distance, response.Routes[0].Duration)
		destination = DestinationRoute{
			destination: coor,
			duration:    response.Routes[0].Duration,
			distance:    response.Routes[0].Distance,
		}
		destinationRoutes = append(destinationRoutes, destination)
	}

	return destinationRoutes, nil
}

func sortDestinations(destinationRoutes []DestinationRoute) []DestinationRoute {
	sort.Slice(destinationRoutes, func(i, j int) bool {
		if destinationRoutes[i].duration == destinationRoutes[j].duration {
			return destinationRoutes[i].destination < destinationRoutes[j].destination
		}
		return destinationRoutes[i].duration < destinationRoutes[j].duration
	})

	return destinationRoutes
}
