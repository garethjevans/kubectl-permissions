build:
	go build -o kubectl-permissions cmd/kubectl-permissions.go

install: build
	sudo cp -f kubectl-permissions /usr/local/bin

lint:
	golangci-lint run

test:
	go test -v -short ./...

.PHONY: integration
integration:
	go test -run Integration ./integration/...
