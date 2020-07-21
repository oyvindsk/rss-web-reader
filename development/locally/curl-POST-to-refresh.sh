
#!/usr/bin/env bash

# try to stop if something fails, 
# see https://stackoverflow.com/questions/821396/aborting-a-shell-script-if-any-command-returns-a-non-zero-value
set -e 
set -o pipefail


AUTH=`echo -n "${RSS_FEED_USERNAME}:${RSS_FEED_PASSWORD}" | base64 -`

# Check for the common pitfall of not having set a username and password: 
if [[ "$AUTH" == "Og==" ]]; then
    echo "Hmm, seems like RSS_FEED_USERNAME and RSS_FEED_PASSWORD are empty, giving up!"
    exit 1
fi

echo "POST localhost:8080/refresh with BasicAuth"
time curl -X POST -H "Authorization: Basic ${AUTH}" localhost:8080/refresh