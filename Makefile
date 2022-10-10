build:
	go build -o permissions cmd/kubectl-permissions.go

lint:
	golangci-lint run
