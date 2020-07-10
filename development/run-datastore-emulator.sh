#!/usr/bin/env bash
gcloud beta emulators datastore start

# To use, just run with the correct env varibale: DATASTORE_EMULATOR_HOST=localhost:8081 go run *.go
