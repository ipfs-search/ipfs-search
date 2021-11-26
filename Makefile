BINARY=build/ipfs-search
SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build -o ${BINARY} main.go

clean:
	rm -f ${BINARY}
	rm -f ${BINARY}.linux64

linux64: ${BINARY}.linux64

$(BINARY).linux64: $(SOURCES)
	env GOOS=linux GOARCH=amd64 go build -o ${BINARY}.linux64 main.go

vagrant: $(BINARY).linux64
	vagrant ssh -c "/vagrant/${BINARY}.linux64 ${ARGS}"
