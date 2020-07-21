#!/usr/bin/env bash

####################################################################################
## RSS config:
####################################################################################
# username and config for the website:
export RSS_FEED_USERNAME="john"
export RSS_FEED_PASSWORD="foobar123"

# file with feeds:
export RSS_FEED_FEEDSFILE="example-feeds.txt"



####################################################################################
## Google Cloud and Docker config:  
####################################################################################
export RSS_FEED_PROJECT=foo-bar-12314
export RSS_FEED_REGION=europe-west1                                      # GCP Region
export RSS_FEED_SERVICE_NAME=web-rss-reader                              # Service Name in Cloud Run
export RSS_FEED_IMAGE_URL="eu.gcr.io/foo-bar-12314/rss-reader:test"      # Docker image tag. Used when building the Docker image locally and when deploying to Cloud Run

# ATM this must be set here, manually, after the first deployment to Cloud Run
# TODO: Figure out how to get this automatically, probably from gcloud run deploy (?)
export RSS_FEED_URL="https://web-rss-reader....a.run.app"
