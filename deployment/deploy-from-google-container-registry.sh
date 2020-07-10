
PROJECT=rss-test-281216                             # GCP project 
REGION=europe-west1                                 # GCP Region
SERVICE_NAME=web-rss-reader                         # Service Name in Cloud Run
IMAGE_URL=eu.gcr.io/rss-test-281216/rss-reader:test # The docker (or..?) image to deploy, uses tagged :latest ? FIXME

gcloud run deploy $SERVICE_NAME --project $PROJECT --platform managed --region $REGION --allow-unauthenticated --image $IMAGE_URL --concurrency 1000