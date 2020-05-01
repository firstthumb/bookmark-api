.PHONY: deploy

build:
	go fmt ./...
	env GOOS=linux go build -ldflags="-s -w" -o bin/bookmark/lambda/main function/lambda/bookmark_lambda.go 
	env GOOS=linux go build -ldflags="-s -w" -o bin/bookmark/worker/main function/worker/bookmark_worker.go 
	env GOOS=linux go build -ldflags="-s -w" -o bin/auth/lambda/main function/authorizer/authorizer.go 
lint:
	golangci-lint run
clean:
	rm -rf ./bin
generate:
	go generate ./...
run:
	gow run cmd/bookmark/main.go
test:
	go test ./...
deploy: clean build
	sls deploy --verbose --force
