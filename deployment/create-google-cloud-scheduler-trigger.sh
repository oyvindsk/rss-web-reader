#!/usr/bin/env bash

# try to stop if something fails, 
# see https://stackoverflow.com/questions/821396/aborting-a-shell-script-if-any-command-returns-a-non-zero-value
set -e 
set -o pipefail


##########################################################

# Arguments -- set these :D
# It's not currently aware of Cloud run as a target, so all thse must be set.
# Most argumenst are sat in the config.sh file, run `source config.sh` or similar ..

# Schedule job name
NAME="rss-refresh"

# Schedule - How often?
# https://crontab.guru/#7_*_*_*_*
SCHEDULE="7 * * * *" # Every hour

# Url, will POST
URL="${RSS_FEED_URL}/refresh"

##########################################################

# https://cloud.google.com/sdk/gcloud/reference/scheduler/jobs/create/http

# Basic auth is just base64 encoding of the username and password
# see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Authorization
# the -n in echo is important :)
AUTH=`echo -n "${RSS_FEED_USERNAME}:${RSS_FEED_PASSWORD}" | base64 -`

echo -e "\nExisting jobs, before:\n"
gcloud --project ${RSS_FEED_PROJECT} scheduler jobs list

echo -e "\nCreating..\n"
gcloud --project ${RSS_FEED_PROJECT} scheduler jobs create http "${NAME}" --schedule="${SCHEDULE}" --uri="${URL}" --attempt-deadline="10m" --description="Refresh the RSS feeds" --headers="Authorization=Basic ${AUTH}"


echo -e "\nExisting jobs, after\n"
gcloud --project ${RSS_FEED_PROJECT} scheduler jobs list
