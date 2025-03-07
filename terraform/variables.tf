variable "region" {
  description = "AWS Region"
  type        = string
  default     = "us-west-2"
  validation {
    condition     = contains(["us-east-1", "us-west-2"], var.region)
    error_message = "Valid regions are us-east-1, us-west-2."
  }
}

locals {
  name        = "os-test"
  lambda_name = "os-test-ingest"
}
