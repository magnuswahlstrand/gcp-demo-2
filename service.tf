data "google_client_config" "default" {}

provider "docker" {
  registry_auth {
    address = "gcr.io"
    username = "oauth2accesstoken"
    password = data.google_client_config.default.access_token
  }
}

data "docker_registry_image" "upload_service_image" {
  name = "gcr.io/${var.project}/${local.service_name}"

}

data "google_container_registry_image" "upload_service_image" {
  name = local.service_name
  digest = data.docker_registry_image.upload_service_image.sha256_digest
}

# The Cloud Run service
resource "google_cloud_run_service" "upload" {
  name = local.service_name
  location = var.region
  autogenerate_revision_name = true

  template {
    spec {
      service_account_name = google_service_account.worker.email
      containers {
        image = data.google_container_registry_image.upload_service_image.image_url
        //        image = data.external.image_digest.result.image
        env {
          name = "BUCKET_NAME"
          value = google_storage_bucket.upload_storage.name
        }
      }
    }
  }
  traffic {
    percent = 100
    latest_revision = true
  }

  depends_on = [
    google_project_service.run]
}

# Set service public
data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.upload.location
  project = google_cloud_run_service.upload.project
  service = google_cloud_run_service.upload.name

  policy_data = data.google_iam_policy.noauth.policy_data
  depends_on = [
    google_cloud_run_service.upload]
}
