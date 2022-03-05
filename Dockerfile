# syntax=docker/dockerfile:1

# Build
FROM golang:1.17.8-buster AS build

WORKDIR /app

COPY src/go.mod ./
COPY src/go.sum ./
RUN /go mod download

COPY src/*.go ./

RUN go build -o /server

# Run
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /server /server

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/server"]
