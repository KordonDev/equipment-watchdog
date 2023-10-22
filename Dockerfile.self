FROM golang:1.19-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/equipment-watchdog

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v

# Remove when using other database than sqlite
RUN apk add build-base

# Build the Go app
RUN go build -o ./out/equipment-watchdog .

# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates


COPY --from=build_base /tmp/equipment-watchdog/out/equipment-watchdog /app/equipment-watchdog

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/go-sample-app"]
