## Task Background

The company is building shipping solutions. One of the problems the company faces is
presenting
a list of pickup locations available in customer area so that he can have his package
delivered to the nearest location in his neighborhood.

## Description

We need to build a small web service that takes the source and a list of destinations
and returns a list of routes between source and each destination. Both source and
destination are defined as a pair of latitude and longitude. The returned list of
routes
should be sorted by driving time and distance (if time is equal).
Thus, we want to answer the question: which destination is closest to the source and
how fast one can get there by car.
The API call should look like this:
GET
http://your-service/routes?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219
The response should look like this:
HTTP/1.1 200 OK
Content-Type: application/json
{
"source": "13.388860,52.517037",
"routes": [
{
"destination": "13.397634,52.529407",
"duration": 465.2,
"distance": 1879.4
},
{
"destination": "13.428555,52.523219",
"duration": 712.6,
"distance": 4123
}
]
}
Where input parameters are:

- src - source location (customer's home), only one can be provided
- dst - destination location (pickup point), multiple can be provided
  For this assignment, you must use http://project-osrm.org/ third-party router service
  to
  get driving times and distances.
  Example request for OSRM you would be using is:
  curl
  'http://router.project-osrm.org/route/v1/driving/13.388860,52.517037;13.397634,52.529407?overview=false'
  Example OSRM response is:
  {
  "routes": [
  {
  "legs": [
  {
  "summary": "",
  "weight": 634,
  "duration": 465.2,
  "steps": [],
  "distance": 1879.4
  }
  ],
  "weight_name": "routability",
  "weight": 634,
  "duration": 465.2,
  "distance": 1879.4
  }
  ],
  "waypoints": [
  {
  "hint":
  "DxQKgArnaoYoAAAAPwAAAA8AAAAAAAAAKAAAAD8AAAAPAAAAAAAAAC7rAAAATMwAqVghAzxMzACtWCEDAQDfC
  p4VrCU=",
  "name": "Friedrichstraße",
  "location": [
  13.3888,
  52.517033
  ]
  },
  {
  "hint":
  "VFQUgPPB9YENAAAACwAAAF0BAAAAAAAADQAAAAsAAABdAQAAAAAAAC7rAAB_bswAGIkhA4JuzAD_iCEDAgCfE
  J4VrCU=",
  "name": "Torstraße",
  "location": [
  13.397631,
  52.529432
  ]
  }
  ],
  "code": "Ok"
  }
  Relevant parts of the response are `routes.0.duration`, `routes.0.distance` and `code`
  for error handling.

## Technology Stack

- Go
- some other 3rd party packages or frameworks 

