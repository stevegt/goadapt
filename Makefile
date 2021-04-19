packages =$(shell go list ./...)
coverpkgs = $(shell echo $(packages) | tr ' ' ',')
pkg = $(shell basename `go list`)
cov = /tmp/$(pkg)-coverage

test: main eg

main:
	go vet
	go test -coverprofile=$(cov).out -coverpkg=$(coverpkgs)
	go tool cover -html=$(cov).out -o $(cov).html

eg:
	for f in examples/*.go; do go vet $$f || exit 1; done 

testloop:
	while true; do inotifywait -e MOVE *; sleep 1; make test; done

