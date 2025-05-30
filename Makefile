TMP := ./.tmp
BINARY_NAME := thop
BINARY_PATH := $(TMP)/bin/$(BINARY_NAME)

# Development

.PHONY: build
build:
	go build -o=$(BINARY_PATH)

# make run [command] <args> to run the app
# this execution is a bit hacky, but it makes it nice to use
.PHONY: run
run: build
	$(BINARY_PATH) $(filter-out run,$(MAKECMDGOALS))

# regular run, with args, for testing more complex commands
.PHONY: run/args
run/args: build
	$(BINARY_PATH) $(ARGS)

.PHONY: test
test: build
	go test -coverprofile=$(TMP)/coverage.out ./...

.PHONY: test/v
test/v: build
	go test -coverprofile=$(TMP)/coverage.out -v ./...

.PHONY: cover
cover: test
	go tool cover -func=$(TMP)/coverage.out

.PHONY: cover/html
cover/html: test
	go tool cover -html=$(TMP)/coverage.out

# Installation

.PHONY: install
install: build
	sudo cp $(BINARY_PATH) /usr/local/bin/
	@echo "Successfully installed $(BINARY_NAME) on your system"

# Quality control

.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-ST1001,-ST1020,-ST1021 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

# capture all errors, for nicer output
%:
	@:
