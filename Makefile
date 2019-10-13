packages =$(shell go list ./...)
coverpkgs = $(shell echo $(packages) | tr ' ' ',')
pkg = $(shell basename `go list`)
cov = /tmp/$(pkg)-coverage

test: 
	go vet
	go test -coverprofile=$(cov).out -coverpkg=$(coverpkgs)
	go tool cover -html=$(cov).out -o $(cov).html

testloop:
	while true; do inotifywait *; sleep 1; make test; done

