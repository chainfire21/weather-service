# Weather Service

This webservice has one endpoint that receives latitude and longitude and returns a short description of the weather there.

## Deployment

Compile the program then run.

On windows, use `go build ./...` then run the executable that gets created.

On mac after compiling you can use `go run ./...`

To hit the endpoint, either `curl` commands or Postman work well, it should be running on `localhost:8080`

eg `curl "http://localhost:8080/weather?lat=38.8894&lng=-77.0352"`

### Endpoint: Get Weather

* Path: `/weather?lat={lat}&lng={lng}`
* Method: `GET`
* Response: string

This endpoint returns a short description of the weather for today

looks up the receipt by the ID and returns an object specifying the points awarded.

Example Response:
`The day will be partly cloudy with a moderate average temperature.`

### Notes:

Usually my project organization would be split up into separate folders for `errors` and `repository` but in the interest of time I threw it all into `service`.

I also skipped adding tests in the interest of time as well.

There are also some more checks I would add, checking to make sure lat and lng are provided, also checking to make sure the temperature is returned in F and not C (which is why I have the temperature unit in my struct) that I did not get to.