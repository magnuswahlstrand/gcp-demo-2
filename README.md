# Google Cloud Platform - Demo 2 - Terraform and signed urls

This demo showcases two cool concepts.

* Configuration & deployment to [Google Cloud Platform](https://cloud.google.com/) using [Terraform](https://www.terraform.io/).
* Using [signed urls](https://cloud.google.com/storage/docs/access-control/signed-urls) for direct file upload and download to [Google Cloud Storage](https://cloud.google.com/storage).

### Inspiration
The demo was heavily inspired by the [Serverless Expeditions](https://www.youtube.com/hashtag/serverlessexpeditions) by Google and [@martinomander](https://twitter.com/martinomander).

Especially these two episodes:
* Three alternatives for running your web app serverless
  ([Youtube](https://youtu.be/ca8FgxpmKVE))
* Terraform, serverless, and Cloud Run in practice ([Youtube](https://youtu.be/IBm0SmwEWpA))

## Deploy using Terraform:

1. Create a new project on Google Cloud with billing enabled
2. Open the GCP Cloud console
3. Clone this repo
    * If you want, you can both open the console and clone using [this link](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/kyeett/gcp-demo-2&page=editor&open_in_editor=README.md).

### Build and deploy docker image for Cloud Run*
```
gcloud builds submit --project your_project_name
```

You will be asked enable the [Cloud Build API](https://cloud.google.com/build) for the project. Press y + [enter]
```
API [cloudbuild.googleapis.com] not enabled on project [...].
Would you like to enable and retry (this will take a few minutes)?
(y/N)?
```

This will take a few minutes


### Deploy with Terraform
```
terraform init
terraform apply
```

## TODO
* [ ] Make `files` cache safe for concurrency
* [ ] URL shortener
* [ ] Limit upload max size
* [ ] Update README 
  * [ ] Complete installation steps
  * [ ] Description of signed URLs
* [ ] Clean up older uploads
* [ ] Fix correct redirect url after POST
* [x] Avoid duplicate names
* [x] Add QR code
