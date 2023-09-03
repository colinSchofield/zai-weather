![](zai-logo.png)

## Zai Backend Code Test

### Synopsis

This document discusses my thought process during the completion of this test, any assumptions that I made, along with a brief summary of the technical challenges and how I overcame them.

The task itself involved the creation of a service to act as a REST API that provides weather information (i.e. temperature in Celsius and wind speed in km/hr).

### Technical Details

This program was written as a Go 1.21 microservice that provides a single Restful endpoint, for serving the weather details. 

The service returns a JSON payload with a unified response as per the specifications. An example is shown below:

```
{
  "status": 200,
  "message": "Request successful",
  "data": {
    "temperature_degrees": 29,
    "wind_speed": 20
  }
}
```

The status and message returns the human readable HTTP Status Code, along with diagnostic information about the payload (i.e. 'Request successful', 'Request successful (cached)', 'Request failure (cache is stale)' or 'Location could not be found')

For further details on this, please refer to the included Open API 3 specification [here](http://github.com/ColinSchofield/zai-weather/open-api/index.html).

Open API can be easily loaded into Stoplight for the API design first approach and then to Sauce Labs (or Postman collection via Newman) for all your API contract testing needs.

### Makefile

As Go does not have any kind of a build automation tool, the UNIX Make provides us with a great starting point.

Execute `make` or `make help` for a list of available tasks. As shown below:

```shell
Usage: make [TARGET]
Targets:
  help                      Show this help message.
  run                       Build and Run (in Docker) the Zai weather service.
  lint                      Run lint checks.
  test                      Test and Code Coverage.
  build                     Build Docker image.
  shell                     Shell into Docker image.
  clean                     Remove any transient build artifacts.
```

### How to run the service

Your machine must have the following installed to be able to run this service:
 - Make
 - Docker

 Then perform the following:
 ```
 1. make run
 ```

This will spin up a docker image that supports Go 1.21. The code will first be checked against lint (i.e. golangci-lint), test cases and code coverage will be run, before the application is built and the image loaded onto your machine. Finally, the application will be started (it should be running on port 8080).

To test it:

`curl -i "http://localhost:8080/v1/weather?city=Melbourne"`

```
HTTP/1.1 200 OK
{"status":200,"message":"Request successful","data":{"temperature_degrees":19,"wind_speed":30}}
```

You may also try the following:

1. `curl -i "http://localhost:8080/v1/weather?city=Sydney"`
2. `curl -i "http://localhost:8080/v1/weather?city=Sydney"` (if within 3 seconds this will be a cached value)
3. `curl -i "http://localhost:8080/v1/weather?city=UnknownPlace"` (returns a 404)

#### Test Cases

During the design and coding of this service, particular attention has been paid to TDD, resulting in the creation of **18** Unit Test cases.

The code coverage was run, revealing that the average code coverage was **91%**.

### Issues & Assumptions

In this section I will highlight a couple of issues that I was considering during the development of this service.

I was mindful of keeping my solution as simple as possible, as per the requirements, without going overboard ([Yagni](https://en.wikipedia.org/wiki/You_aren%27t_gonna_need_it)).

1. The program was developed using Idiomatic Go and it was kept clean, clear and as simple as possible with comments provided when deemed appropriate.

2. Where appropriate configuration values have been used (such as the TTL of the cache, the web access keys etc). These  were specified as environment variables, as a convenience to Docker integration.

3. Docker was used during development. This was to make the service more easily includable in a production environment (i.e. We may want to load it onto EKS, ECS, Fargate, Lambda etc)

4. The task involved two 3rd party weather providers (a primary and a fail-over). I considered the scenario where either the primary or *both* the primary and fail-over went down. This would potentially lead to a double timeout of over six seconds per request(!) To mitigate against this, the [circuit breaker design pattern](https://en.wikipedia.org/wiki/Circuit_breaker_design_pattern) was employed.

5. Whenever possible go libraries have been used -- [Gin](https://gin-gonic.com) for the HTTP Web framework, [Resty](https://dev.to/ankitmalikg/go-how-to-use-resty-2pmg) for the REST client, [go-cache](https://github.com/patrickmn/go-cache) for the software caching and [gobreaker](https://dev.to/he110/circuitbreaker-pattern-in-go-43cn) for the Circuit Breaker.

6. If the service requires CORS (Cross-Origin Resource Sharing) access, then this could be provided for in Gin with a Middleware HandlerFunc.

### Scalability and Reliability

1. The scalability of the service could be improved by adding horizontal scaling behind a load balancer and through fine tuning the cache TTL.

To gain *more* performance, would require a design discussion, but here is what I am thinking. Each request should **not** require a callout to a 3rd party provider. Instead, just before the cache is going to expire, a request is made that gathers **all** the cities in Australia and store this information. This would provide instant access to our callers.

Alternatively, rather than *polling* our providers, it might be possible to subscribe to a weather broker and receive notifications when a weather change has occurred. Then the *latest* information could be updated from the weather providers.

2. In terms of reliability, I feel that we would definitely benefit from the addition of [Prometheus](https://grafana.com/go/webinar/intro-to-observability-with-prometheus/) metrics for monitoring and observability. Looking at the information recorded by the Circuit Breaker component [gobreaker.Counts](https://github.com/sony/gobreaker/blob/70f7cbc53af96e27e1042a5f5803c9b960e0ca81/gobreaker.go#L47). This would be excellent for this purpose and could be used to build up alerts and/or dashboards (via Grafana).



Thank you guys for giving me this coding challenge, I had a lot of fun working on it. 🙂

**Colin Schofield**
e: colin_sch@yahoo.com
p: 0448 644 233
l: https://www.linkedin.com/in/colins/
