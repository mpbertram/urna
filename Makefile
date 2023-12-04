.PHONY: install
install:
	go install github.com/google/certificate-transparency-go

.PHONY: build
build: install
	go build ./...

.PHONY: test
test: install
	go test -v ./... -timeout 30m

.PHONY: coverage
COVERPROFILE ?= coverage.out
coverage: install
	go test -v ./... -timeout 30m -cover -coverprofile=${COVERPROFILE}
