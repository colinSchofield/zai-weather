FROM golang:1.21.0

ENV PORT :8080
ENV PRIMARY_TIMEOUT_SECONDS 3
ENV PRIMARY_ACCESS_KEY 1cadfad44c3387c66d14a12cb33f282e
ENV PRIMARY_END_POINT http://api.weatherstack.com/current

ENV FAILOVER_TIMEOUT_SECONDS 3
ENV FAILOVER_ACCESS_KEY fe0e197efcdefea9a19e9c4810f2801b
ENV FAILOVER_END_POINT http://api.openweathermap.org/data/2.5/weather

ENV PRIMARY_REQUESTS 3
ENV PRIMARY_FAILURE_RATIO 0.6
ENV FAILOVER_REQUESTS 3
ENV FAILOVER_FAILURE_RATIO 0.6

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN golangci-lint run ./...
RUN go test ./... -cover
RUN go build -o main src/main.go
CMD ["/app/main"]