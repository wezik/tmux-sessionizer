# runs the app
run *args:
  just build
  ./build/phop {{args}}

# builds the app
build:
  go build -o ./build/phop

# installs the app on the system
install:
  just build
  sudo cp ./build/phop /usr/local/bin/

# formats all files with go formatter
fmt:
  go fmt ./...
