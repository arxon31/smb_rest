start:
	cd cmd/server && ./server

build:
	go build -C cmd/server/ -o server

docker:
	docker build -t smbrest:local .

test:
	go test -v ./...

lint:
	golangci-lint run ./... --timeout 5m

run: test build docker start
