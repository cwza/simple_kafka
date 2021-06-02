# build stage
FROM golang:1.16-alpine AS build-env
ADD . /src
RUN cd /src && go mod download
RUN cd /src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./producer/producer.go -o producer
RUN cd /src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./consumer/consumer.go -o consumer

# final stage
FROM alpine
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /src/producer/producer /
COPY --from=build-env /src/consumer/consumer /