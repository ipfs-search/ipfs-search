FROM golang:1.16-alpine AS build

RUN apk add --no-cache git gcc musl-dev

WORKDIR /src/

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod go.sum /src/

#This is the ‘magic’ step that will download all the dependencies that are specified in
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction this line)
RUN go mod download && go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get -v

# Here we copy the rest of the source code
COPY . /src/

# Run the build
RUN go install -v ./...

# This results in a single layer image
FROM alpine
COPY --from=build /go/bin/ipfs-search /usr/local/bin/ipfs-search

CMD ["crawl"]
ENTRYPOINT ["/usr/local/bin/ipfs-search"]
