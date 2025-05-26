# runs the app
run *args:
  just build
  ./.build/phop {{args}}

# builds the app
build:
  go build -o ./.build/phop .

# installs the app on the system
install:
  just build
  sudo cp ./.build/phop /usr/local/bin/

# formats all files with go formatter
fmt:
  go fmt ./...

# runs test suite
test *args:
  go test ./... {{args}}

test-coverage:
  just test-coverage-func

test-coverage-func:
  go test -coverprofile=coverage.out -v ./...
  go tool cover -func=coverage.out
  rm coverage.out

test-coverage-html:
  go test -coverprofile=coverage.out -v ./...
  go tool cover -html=coverage.out
  rm coverage.out

coverage:
  go test -coverprofile=coverage.out ./... > /dev/null
  go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}'
  rm coverage.out
