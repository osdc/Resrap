.PHONY: build test lint vet fmt coverage cover-html clean

build:
	go build ./...

test:
	go test ./... -count=1

lint:
	golangci-lint run ./...

vet:
	go vet ./...

fmt:
	gofmt -l -w .

coverage:
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -func=coverage.out

cover-html: coverage
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -f coverage.out coverage.html
