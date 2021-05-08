
# Create a service account
resource "google_service_account" "worker" {
  account_id   = "upload-service"
  display_name = "Upload Service SA"
}

resource "google_service_account_key" "worker_key" {
  service_account_id = google_service_account.worker.id
}

resource "google_secret_manager_secret" "worker_key" {
  secret_id = "${google_service_account.worker.account_id}_key"

  replication {
    user_managed {
      replicas {
        location = var.region
      }
    }
  }

  depends_on = [
    google_service_account.worker]
}

resource "google_secret_manager_secret_version" "secret_version_one" {
  secret = google_secret_manager_secret.worker_key.id

  secret_data = base64decode(google_service_account_key.worker_key.private_key)
}

resource "google_secret_manager_secret_iam_member" "member" {
  project = google_secret_manager_secret.worker_key.project
  secret_id = google_secret_manager_secret.worker_key.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.worker.email}"
}