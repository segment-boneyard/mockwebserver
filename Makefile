test:
	@godep go test -cover ./...

build:
	@godep go build

.PHONY: test build