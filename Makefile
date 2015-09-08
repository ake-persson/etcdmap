all: test readme

format:
	gofmt -w=true .

syntax: format
	golint .
	go vet .
	go build

test: syntax
	go test

readme:
	cp README_HEAD.md README.md
	godoc2md github.com/mickep76/etcdmap | grep -v Generated >>README.md
