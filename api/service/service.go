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

type Parameters struct {
	Src []string
	Dst []string
}

type destinationRoutes struct {
	fetchedDestinationRoutes []DestinationRoute
}

func (p Parameters) ValidateParameters() error {
	if p.Src == nil {
		return errors.New("source location is empty")
	}
	if len(p.Src) > 1 {
		return errors.New("only one source location can be provided")
	}
	_, _, err := coordsparser.Parse(p.Src[0])
	if err != nil {
		return err
	}

	if p.Dst == nil {
		return errors.New("destination location is empty")
	}

	for _, dst := range p.Dst {
		_, _, err := coordsparser.Parse(dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Parameters) GetList() (RouteResponse, error) {
	fetchedDestinationRoutes, err := p.getDestinationRoutes()
	if err != nil {
		return RouteResponse{}, err
	}
	destinationRoutes{fetchedDestinationRoutes}.sortDestinations()

	return RouteResponse{Source: p.Src[0], Routes: fetchedDestinationRoutes}, nil
}

func (p Parameters) getDestinationRoutes() ([]DestinationRoute, error) {
	var thirdPartyResponse ThirdPartyResponse
	var destinationRoutes []DestinationRoute
	var destinationRoute DestinationRoute

	for _, coor := range p.Dst {
		url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%s;%s?overview=false", p.Src[0], coor)
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
		destinationRoute = DestinationRoute{
			Destination: coor,
			Duration:    thirdPartyResponse.Routes[0].Duration,
			Distance:    thirdPartyResponse.Routes[0].Distance,
		}
		destinationRoutes = append(destinationRoutes, destinationRoute)
	}

	return destinationRoutes, nil
}

func (dr destinationRoutes) sortDestinations() {
	sort.Slice(dr.fetchedDestinationRoutes, func(i, j int) bool {
		if dr.fetchedDestinationRoutes[i].Duration == dr.fetchedDestinationRoutes[j].Duration {
			return dr.fetchedDestinationRoutes[i].Destination < dr.fetchedDestinationRoutes[j].Destination
		}
		return dr.fetchedDestinationRoutes[i].Duration < dr.fetchedDestinationRoutes[j].Duration
	})
}
