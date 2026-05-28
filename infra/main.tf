terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }

  backend "s3" {
    bucket   = "difference-engine-tfstate"
    key      = "terraform.tfstate"
    region   = "auto"

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

provider "cloudflare" {
  # Set CLOUDFLARE_API_TOKEN in your environment
}

# ── R2 Buckets ────────────────────────────────────────────────────────────────

resource "cloudflare_r2_bucket" "stems" {
  account_id = var.account_id
  name       = "difference-engine-stems"
}

resource "cloudflare_r2_bucket" "output" {
  account_id = var.account_id
  name       = "difference-engine-output"
}

# ── R2 API Token (used by the Workers container) ──────────────────────────────

data "cloudflare_api_token_permission_groups" "all" {}

resource "cloudflare_api_token" "r2_worker" {
  name = "difference-engine-worker-r2"

  policy {
    permission_groups = [
      data.cloudflare_api_token_permission_groups.all.r2["Workers R2 Storage Bucket Item Read"],
      data.cloudflare_api_token_permission_groups.all.r2["Workers R2 Storage Bucket Item Write"],
    ]
    resources = {
      "com.cloudflare.edge.r2.bucket.${var.account_id}_default_${cloudflare_r2_bucket.stems.name}"  = "read"
      "com.cloudflare.edge.r2.bucket.${var.account_id}_default_${cloudflare_r2_bucket.output.name}" = "edit"
    }
  }
}

# ── Pages (static frontend) ───────────────────────────────────────────────────

locals {
  workers_url = "https://difference-engine-api.${var.workers_subdomain}.workers.dev"
}

resource "cloudflare_pages_project" "frontend" {
  account_id        = var.account_id
  name              = "difference-engine"
  production_branch = "main"

  deployment_configs {
    production {
      compatibility_date = "2024-01-01"
      environment_variables = {
        WORKERS_URL = local.workers_url
      }
    }
  }
}

resource "cloudflare_pages_domain" "frontend" {
  account_id   = var.account_id
  project_name = cloudflare_pages_project.frontend.name
  domain       = var.domain
}

# ── Workers (container — deployed via wrangler, registered here) ──────────────

resource "cloudflare_workers_script" "api" {
  account_id = var.account_id
  name       = "difference-engine-api"
  content    = "export default { async fetch(req) { return fetch(req); } }"
  module     = true
}
