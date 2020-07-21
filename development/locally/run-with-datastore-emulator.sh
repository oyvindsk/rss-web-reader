#!/usr/bin/env bash

# try to stop if something fails, 
# see https://stackoverflow.com/questions/821396/aborting-a-shell-script-if-any-command-returns-a-non-zero-value
set -e 
set -o pipefail

# cleanup on exit
function on-exit {
    echo -e "\n$0 exiting, removing tmp-rss"
    rm ./tmp-rss
}
trap on-exit EXIT


cd src
go build -o ../tmp-rss

cd ..
echo -e "\n\nRunning with: user/pass/feeds-file: $RSS_FEED_USERNAME / $RSS_FEED_PASSWORD / $RSS_FEED_FEEDSFILE \n\n"
echo -e "\n!!!!!!!!!!!!!!!!!!!!!\n Running with local Datastore emulator. Make sure you started it with ./development/run-datastore-emulator.sh =)\n\n!!!!!!!!!!!!!!!!!!!!!\n\n\n\n\n\n"
DATASTORE_EMULATOR_HOST=localhost:8081 ./tmp-rss