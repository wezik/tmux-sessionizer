# runs the app
run *args:
  just build
  ./build/phop {{args}}

# builds the app
build:
  go build -o build/phop ./src/

# installs the app on the system
install:
  just build
  sudo cp ./build/phop /usr/local/bin/

# formats all files with go formatter
fmt:
  go fmt ./...

# runs test suite
test *args:
  go test ./test/... {{args}}
