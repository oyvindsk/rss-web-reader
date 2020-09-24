#!/usr/bin/env bash

# try to stop if something fails, 
# see https://stackoverflow.com/questions/821396/aborting-a-shell-script-if-any-command-returns-a-non-zero-value
set -e 
set -o pipefail

gcloud --project rss-test-281216 datastore indexes create index.yaml

echo -e "\n\nCurrent indexes, to cleanup run:\ngcloud --project rss-test-281216 datastore indexes cleanup index.yaml\n\n"

gcloud --project rss-test-281216 datastore indexes list
