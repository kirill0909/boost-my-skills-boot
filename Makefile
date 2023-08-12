all: fmt goimports golangci tidy

fmt:
	@gofmt -w -s app
goimports:
	@goimports -w app
tidy:
	@cd app && go mod tidy
golangci:
	cd app && golangci-lint run -v