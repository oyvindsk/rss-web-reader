#!/usr/bin/env bash

# try to stop if something fails, 
# see https://stackoverflow.com/questions/821396/aborting-a-shell-script-if-any-command-returns-a-non-zero-value
set -e 
set -o pipefail

gcloud beta emulators datastore start

# To use, just run with the correct env varibale: DATASTORE_EMULATOR_HOST=localhost:8081 go run *.go
