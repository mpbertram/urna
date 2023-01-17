.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test -v ./urna/*.go -timeout 30m

.PHONY: coverage
COVERPROFILE ?= coverage.out
coverage:
	go test -v ./urna/*.go -timeout 30m -cover -coverprofile=${COVERPROFILE}
