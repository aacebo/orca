clean:
	rm -rf ./bin

build:
	go build ./...

clean_build: clean build

fmt:
	gofmt -w ./

test:
	go clean -testcache
	go test ./... -cover

run:
	go run ./
