default: fmt lint install generate

build:
	@go build -v ./...

install: build
	@go install -v ./...

lint:
	@golangci-lint run

generate:
	@cd tools; go generate ./...

fmt:
	@gofmt -s -w -e .

test:
	@go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	@TF_ACC=1 go test -v -cover -timeout 120m ./...

hooks:
	@mkdir -p .git/hooks
	@echo '#!/bin/sh' > .git/hooks/pre-commit
	@echo 'set -e' >> .git/hooks/pre-commit
	@echo 'make generate' >> .git/hooks/pre-commit
	@echo 'make fmt' >> .git/hooks/pre-commit
	@echo 'git add -u' >> .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit

vulncheck:
	@govulncheck -show verbose ./...

.PHONY: fmt lint test testacc build install generate
