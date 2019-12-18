all: clean build test

build:
	go build ./...
	cd cmd/tokengen; go build
	cd cmd/server; go build

run: build
	cd cmd/server; go build; ./server 
	
test: build
	go test -v 

coverage:
	go test ./... -coverprofile=coverage.out | true
	go tool cover -html=coverage.out

clean:
	cd cmd/tokengen; rm -f tokengen
	cd cmd/server; rm -f server 
	cd _api; make clean
	rm -f coverage.out

