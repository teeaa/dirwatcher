.DEFAULT_GOAL := generate

GOBIN := $(shell pwd)

generate:
	@echo "Building binaries..."
	@GOBIN=$(GOBIN) go get ./...
	@GOBIN=$(GOBIN) go build ./...
	@echo "Done!"

clean:
	@echo "Cleaning up..."
	@rm $(GOBIN)/dirwatcher
	@rm $(GOBIN)/dirserver
	@echo "Done!"

install:
	@echo "If you really want a global install run 'make real_install'"

real_install:
	@echo "Installing..."
	@go get ./...
	@go install ./...
	@echo "Done!"