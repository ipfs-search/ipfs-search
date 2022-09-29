# syntax=docker/dockerfile:1.3
FROM golang:1.19-alpine AS build

RUN apk add --no-cache git gcc musl-dev

WORKDIR /src/

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod go.sum /src/

RUN go mod download && go mod verify

# Here we copy the rest of the source code
COPY . /src/

# Run the build
RUN --mount=type=cache,target=/go-build GOCACHE=/go-build GORACE="halt_on_error=1" go install -v -race ./...

# This results in a single layer image
FROM alpine AS runtime
COPY --from=build /go/bin/ipfs-search /usr/local/bin/ipfs-search

CMD ["crawl"]
ENTRYPOINT ["/usr/local/bin/ipfs-search"]
