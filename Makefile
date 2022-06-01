GOCMD=$(shell echo go)
GOLINT=$(shell echo golangci-lint)

fmt:
	@echo "+ $@"
	@$(GOCMD) fmt ./...

lint: 
	@echo "+ $@"
	@${GOLINT} run

test:
	@echo "+ $@"
	@$(GOCMD) test ./... -race -v -coverprofile=coverage.out -covermode=atomic

benchmark:
	@echo "+ $@"
	@$(GOCMD) test ./... -bench=. -run=^#

build:
	@echo "+ $@"
	@$(GOCMD) build

all: fmt lint test benchmark build	

demo:
	@echo "+ $@"
	@$(GOCMD) run cmd/demo.go
