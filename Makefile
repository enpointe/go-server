all: clean build tokengen server test

build:
	go build ./...

tokengen: build
	cd cmd/tokengen; go build

server: build
	cd cmd/server; go build

run: server
	cd cmd/server; go build; ./server 
	
test: build
	go test -v 

coverage:
	go test ./... -coverprofile=coverage.out | true
	go tool cover -html=coverage.out

docker:
	docker build --file Dockerfile -t go-server-web .
	docker build --file Dockerfile.tokengen -t go-server-tokengen .
	@echo "Successfully create docker repositories for 'go-sever-web' and 'go-server-tokengen'"
	@docker images go-server-web
	@docker images go-server-tokengen

clean:
	cd cmd/tokengen; rm -f tokengen
	cd cmd/server; rm -f server 
	cd _api; make clean
	rm -f coverage.out
