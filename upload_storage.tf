resource "google_storage_bucket" "upload_storage" {
  name = local.bucket_name

  location = var.location
}

resource "google_storage_bucket_iam_member" "member" {
  bucket = google_storage_bucket.upload_storage.name
  role = "roles/storage.admin"

  member = "serviceAccount:${google_service_account.worker.email}"
}