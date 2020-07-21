
exit 1
# FIXME


# # Move the secret Dockerfile to the standard filename for gcloud
# mv SECRET-Dockerfile Dockerfile

# # Upload the files from this diretory and build it remotly in Google Cloud Build, 
# # see https://console.cloud.google.com/cloud-build/builds 
# # or    gcloud builds list
# #       gcloud builds list --ongoing
# #
# # If it's successful, the image is pushed to Google Container Registry: https://console.cloud.google.com/gcr/images/
# #
# # Will ignore (leave out) files from .gitignore by default!
# gcloud builds submit --project $PROJECT --tag $IMAGE_URL

# # Move the Dockerfile back
# mv Dockerfile SECRET-Dockerfile