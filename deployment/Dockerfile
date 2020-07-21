# Use the offical Golang image both build AND run it
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang

FROM golang:latest

# Arguments uses in this Dockerfile, must be passed when `docker build` with `--build-arg`
ARG RSS_FEED_FEEDSFILE

# Copy local code to the container image.
# Just copy what we need: 
# - *.go        : the code
# - *.html      : templates, static files etc
# - go.mod|sum  : so go build understands what module we are in 
#                 (can build the pkg from "import path" but still use the local files on disk)

WORKDIR /app

COPY ./src/go.mod     .
COPY ./src/go.sum     .

RUN go get

COPY ./src/*.go            ./
COPY ./templates/*.html    ./templates/
COPY ./$RSS_FEED_FEEDSFILE .

# Build the command inside the container.
# Build with import path, maybe we should just give the local files instead?    
# go does not fetch the package from github since we are already i the module 
RUN CGO_ENABLED=0 GOOS=linux go build -o runme "foo.bar/rss"

# Module / pkg debug: RUN go list -f '{{.Module.Dir}}' "github.com/oyvindsk/web-oyvindsk.com/blogbackend"

# Set a few environment variables
# could be secret if we don't have any "real" handling of secrets

# Run the compiled binary on container startup.
# again, it expects to be started from the same dir as all it files it serves
WORKDIR /app

# Environment varaiables that should be set at container runtime, usually by docker run --env .. :
# or, for Cloud Run: gcloud run deploy --set-env-vars=[KEY=VALUE,...]
#   RSS_FEED_PROJECT
#   RSS_FEED_FEEDSFILE
#   RSS_FEED_USERNAME
#   RSS_FEED_PASSWORD

CMD ["./runme"]
