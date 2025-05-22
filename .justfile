# runs the app
run *args:
  ./build/phop {{args}}

# builds the app
build:
  go build -o ./build/phop

# installs the app on the system
install:
  go build -o phop
  sudo cp ./build/phop /usr/local/bin/

# formats all files with go formatter
fmt:
  go fmt ./...
