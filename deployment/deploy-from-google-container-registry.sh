
#!/usr/bin/env bash

# try to stop if something fails, 
# see https://stackoverflow.com/questions/821396/aborting-a-shell-script-if-any-command-returns-a-non-zero-value
set -e 
set -o pipefail

gcloud run deploy               \
    $RSS_FEED_SERVICE_NAME      \
    --set-env-vars="RSS_FEED_PROJECT=${RSS_FEED_PROJECT},RSS_FEED_FEEDSFILE=${RSS_FEED_FEEDSFILE},RSS_FEED_USERNAME=${RSS_FEED_USERNAME},RSS_FEED_PASSWORD=${RSS_FEED_PASSWORD}" \
    --project $RSS_FEED_PROJECT \
    --platform managed          \
    --region $RSS_FEED_REGION   \
    --allow-unauthenticated     \
    --image $RSS_FEED_IMAGE_URL \
    --concurrency 1000
