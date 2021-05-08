output "SERVICE_URL" {
  value = google_cloud_run_service.upload.status[0].url
}

output "BUCKET_NAME" {
  value = google_storage_bucket.upload_storage.name
}

output "SERVICE_ACCOUNT_JSON" {
  value = google_secret_manager_secret_version.secret_version_one.name
}