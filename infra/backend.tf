terraform {
  backend "s3" {
    bucket = "difference-engine-tfstate"
    key    = "terraform.tfstate"
    region = "auto"

    endpoints = {
      s3 = "https://62aa79ca5a4eb69594dcd5b96f00b4bd.r2.cloudflarestorage.com"
    }

    # R2 doesn't use these AWS-specific checks
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
  }
}
