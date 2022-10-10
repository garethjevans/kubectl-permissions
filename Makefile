build:
	go build -o kubectl-permissions cmd/kubectl-permissions.go

install: build
	sudo mv kubectl-permissions /usr/local/bin

lint:
	golangci-lint run
