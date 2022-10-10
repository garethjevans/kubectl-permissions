build:
	go build -o permissions main.go

lint:
	golangci-lint run
