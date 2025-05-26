MAIN_PACKAGE_PATH := "."
BINARY_NAME := "phop"

## Development

run *args:
  @just build
  ./.tmp/bin/{{BINARY_NAME}} {{args}}

# builds the app
build:
  go build -o=./.tmp/bin/{{BINARY_NAME}} {{MAIN_PACKAGE_PATH}}

# installs the app on the system (linux only)
install:
  @just build
  sudo cp ./.tmp/bin/{{BINARY_NAME}} /usr/local/bin/

# formats all files with go formatter
fmt:
  go fmt ./...

# runs test suite
test *args:
  go test {{args}} ./...

# runs test suite with coverage
test-cover:
  go test -coverprofile=./.tmp/coverage.out -v ./...
  go tool cover -html=./.tmp/coverage.out

## Quality control

# checks if there are any uncommitted changes
no-dirty:
  git diff --exit-code

# runs static analysis
audit:
  go mod verify
  go vet ./...
  go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-ST1001,-ST1020,-ST1021 ./...
  go run golang.org/x/vuln/cmd/govulncheck@latest ./...
  go test -vet=off ./...

# tidys up the code
tidy:
  go fmt ./...
  go mod tidy -v

# runs audit, tidy and no-dirty checks before pushing
push:
  @just audit
  @just tidy
  @just no-dirty
  git push
