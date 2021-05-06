terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
    }
  }
}

provider "google" {
  project = var.project
}

locals {
//  function_folder = "function"
//  function_name   = "analyse"

  service_folder = "service"
  service_name   = "upload"

  bucket_folder = "upload"
  bucket_name   = "${var.project}-upload"

//  deployment_name = "cats"
  service_account  = "serviceAccount:${google_service_account.worker.email}"
}
