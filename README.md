# Explanation
## run in docker
under project root directory, <br>
1, enter command to build the image: `docker build --tag app-ingrid .`<br>
2, enter command to start the container: `docker run -it -p 8082:8080 app-ingrid`
## run on host machine
Run program under directory `cmd`(also the main package), run the command `go run shippingWeb.go`.
The HTTP server serves the current working directory on port `8082`, and it can be changed to other port.

## http entry point
to see the response, on your local browser, visit url: http://localhost:8082/list?src=13.5878508,42.527337&dst=13.42855587,52.523219&dst=13.42885587,52.423219&dst=13.428555,52.523421
parameter `src` is a coordinate pair, and it should only be one pair, if not the response will be an error message. `dst` could contain multiple coordinate pairs.

the correct response should look like this:
Content-Type: application/json
Status Code: 200 OK

{
  "source": "13.5878508,42.527337",
  "routes": [
    {
      "destination": "13.42885587,52.423219",
      "duration": 57139.9,
      "distance": 1514304.1
    },
    {
      "destination": "13.428555,52.523421",
      "duration": 57772,
      "distance": 1519757.6
    },
    {
      "destination": "13.42855587,52.523219",
      "duration": 57772,
      "distance": 1519758.2
    }
  ]
}

## run testing

under `testing` folder, run command: `go test -v`