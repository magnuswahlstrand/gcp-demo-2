terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }

    docker = {
      source = "kreuzwerker/docker"
    }
  }
}

provider "google" {
  project = var.project
}

locals {
  service_folder = "service"
  service_name   = "upload"

  bucket_folder = "upload"
  bucket_name   = "${var.project}-upload"

  service_account  = "serviceAccount:${google_service_account.worker.email}"
}
