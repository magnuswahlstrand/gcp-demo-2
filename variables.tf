variable "project" {
  type        = string
  description = "Google Cloud Project ID"
}

variable "region" {
  default = "europe-west1"
  type    = string
}

variable "location" {
  default = "EU"
  type    = string
}
