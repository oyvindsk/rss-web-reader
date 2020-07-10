#!/usr/bin/env bash

# try to stop if something fails, 
# see https://stackoverflow.com/questions/821396/aborting-a-shell-script-if-any-command-returns-a-non-zero-value
set -e 
set -o pipefail

IMAGE_URL=eu.gcr.io/rss-test-281216/rss-reader:test

# remove :latest so we don't acidentally run an old one
# don't use :latest when deploying, keep it ..
# sudo docker rmi ${IMAGE_URL}

sudo docker build -t ${IMAGE_URL} -f ./deployment/SECRET-Dockerfile .

sudo docker run -ti --rm -p 8080:8080 ${IMAGE_URL}
