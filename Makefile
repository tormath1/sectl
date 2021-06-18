BIN := sectl

.PHONY: build
build:
	@go build \
	-o $(BIN)

.PHONY: fmt
fmt:
	@find . -type f -name "*.go" -exec goimports -w -v {} +;
