.PHONY: build clean deploy

build:
	go fmt ./...
	golangci-lint run
	env GOOS=linux go build -ldflags="-s -w" -o bin/bookmark/lambda/main function/lambda/bookmark_lambda.go 
	env GOOS=linux go build -ldflags="-s -w" -o bin/bookmark/worker/main function/worker/bookmark_worker.go 
lint:
	golangci-lint run
clean:
	rm -rf ./bin
generate:
	go generate ./...
run:
	go run cmd/bookmark/main.go
test:
	go test
deploy: clean build
	sls deploy --verbose --force
