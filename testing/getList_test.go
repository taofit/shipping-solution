package testing

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/taofit/shipping-solution/api/service"
)

type ListTest struct {
	src            []string
	dst            []string
	expectedResult service.RouteResponse
}

func TestGetList(t *testing.T) {
	var testVals = []ListTest{
		{
			src: []string{"13.388860,52.517037"},
			dst: []string{"13.42855587,52.523219", "13.42885587,52.423219", "13.428555,52.523421"},
			expectedResult: service.RouteResponse{
				Source: "13.388860,52.517037",
				Routes: []service.DestinationRoute{
					{
						Destination: "13.428555,52.523421",
						Duration:    389.2,
						Distance:    3803.5,
					},
					{
						Destination: "13.42855587,52.523219",
						Duration:    389.2,
						Distance:    3804.2,
					},
					{
						Destination: "13.42885587,52.423219",
						Duration:    1378.1,
						Distance:    14282.9,
					},
				},
			},
		},
		{
			src: []string{"14.388860,53.517037"},
			dst: []string{"14.42855587,53.523219", "14.42885587,53.423219", "14.428555,53.523421"},
			expectedResult: service.RouteResponse{
				Source: "14.388860,53.517037",
				Routes: []service.DestinationRoute{
					{
						Destination: "14.42885587,53.423219",
						Duration:    1814.3,
						Distance:    16243.1,
					},
					{
						Destination: "14.42855587,53.523219",
						Duration:    1899.7,
						Distance:    13066.7,
					},
					{
						Destination: "14.428555,53.523421",
						Duration:    1904.3,
						Distance:    13085.8,
					},
				},
			},
		},
		{
			src:            []string{"12.388860,43.517037"},
			dst:            []string{"14.42255587s,53.523219", "14.62885587,53.423219", "14.428555,53.523421"},
			expectedResult: service.RouteResponse{},
		},
	}

	for i, tt := range testVals {
		testname := fmt.Sprintf("test #%v", i)
		t.Run(testname, func(t *testing.T) {
			routeResponse, _ := service.Parameters{Src: tt.src, Dst: tt.dst}.GetList()
			if !reflect.DeepEqual(routeResponse, tt.expectedResult) {
				t.Errorf("got %v, want %v", routeResponse, tt.expectedResult)
			}
		})
	}
}
