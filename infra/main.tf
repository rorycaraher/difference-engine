terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
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

resource "cloudflare_pages_project" "frontend" {
  account_id        = var.account_id
  name              = "difference-engine"
  production_branch = "main"

  deployment_configs {
    production {
      compatibility_date = "2024-01-01"
    }
  }
}

resource "cloudflare_pages_domain" "frontend" {
  account_id   = var.account_id
  project_name = cloudflare_pages_project.frontend.name
  domain       = var.domain
}

# ── DNS ───────────────────────────────────────────────────────────────────────

resource "cloudflare_record" "pages" {
  zone_id = var.zone_id
  name    = "@"
  value   = "${cloudflare_pages_project.frontend.name}.pages.dev"
  type    = "CNAME"
  proxied = true
}

# ── Workers (container — deployed via wrangler, registered here for routing) ──

resource "cloudflare_workers_script" "api" {
  account_id = var.account_id
  name       = "difference-engine-api"
  content    = "export default { async fetch(req) { return fetch(req); } }"
  module     = true
}

resource "cloudflare_workers_route" "mixdown" {
  zone_id     = var.zone_id
  pattern     = "${var.domain}/mixdown"
  script_name = cloudflare_workers_script.api.name
}
