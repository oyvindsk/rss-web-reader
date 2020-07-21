#!/usr/bin/env bash

# try to stop if something fails, 
# see https://stackoverflow.com/questions/821396/aborting-a-shell-script-if-any-command-returns-a-non-zero-value
set -e 
set -o pipefail


# remove :latest so we don't acidentally run an old one
# don't use :latest when deploying, keep it ..
# sudo docker rmi ${RSS_FEED_IMAGE_URL}

echo "Docker building image withg tag: $RSS_FEED_IMAGE_URL"

sudo docker build --build-arg RSS_FEED_FEEDSFILE=${RSS_FEED_FEEDSFILE} -t ${RSS_FEED_IMAGE_URL} -f ./deployment/Dockerfile .