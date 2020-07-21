#!/usr/bin/env bash

# try to stop if something fails, 
# see https://stackoverflow.com/questions/821396/aborting-a-shell-script-if-any-command-returns-a-non-zero-value
set -e 
set -o pipefail


# remove :latest so we don't acidentally run an old one
# don't use :latest when deploying, keep it ..
# sudo docker rmi ${RSS_FEED_IMAGE_URL}

./deployment/build-docker-image-locally.sh


echo -e "\n\nNot running the Docker container cause it won't work with datastore ATM, see $0"
echo -e "(don't really care since I don't use docker for local go development)\n"

# Ideas to make it use in local Docker:
#  - ADD in the needed auth/authorize files from gcloud
#  - Use the local emulator, should work if we set up the networking in Docker 
#           --env DATASTORE_EMULATOR_HOST=localhost:8081 \

# Call docker run with environment variables that should be set already (probably by: `source config.sh` or similar)
# sudo docker run  \
#     -ti          \
#     -p 8080:8080 \
#     --env RSS_FEED_FEEDSFILE=$RSS_FEED_FEEDSFILE \
#     --env RSS_FEED_USERNAME=$RSS_FEED_USERNAME   \
#     --env RSS_FEED_PASSWORD=$RSS_FEED_PASSWORD   \
#     ${RSS_FEED_IMAGE_URL}