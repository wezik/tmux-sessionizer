#!/bin/bash

# util script that installs the app on the system,
# dont run with sudo to avoid installing go packages system wide

go build -o phop
sudo mv phop /usr/local/bin/
