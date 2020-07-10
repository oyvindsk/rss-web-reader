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
echo -e "\n\nRunning with user/pass: foo bar\n\n"
# RSS_FEED_FEEDSFILE="SECRET-feeds.txt" RSS_FEED_USERNAME=foo RSS_FEED_PASSWORD=bar ./tmp-rss

# Or with datastore emulator: 
# echo -e "Running with local Datastore emulator. Make sure you started it with ./development/run-datastore-emulator.sh =)\n\n"
RSS_FEED_FEEDSFILE="SECRET-feeds.txt" RSS_FEED_USERNAME=foo RSS_FEED_PASSWORD=bar DATASTORE_EMULATOR_HOST=localhost:8081 ./tmp-rss
