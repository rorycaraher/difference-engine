# ── DNS ───────────────────────────────────────────────────────────────────────

resource "cloudflare_record" "pages" {
  zone_id = var.zone_id
  name    = var.domain
  content = "${cloudflare_pages_project.frontend.name}.pages.dev"
  type    = "CNAME"
  proxied = true
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

# ── Pages (static frontend) ───────────────────────────────────────────────────

locals {
  workers_url = "https://difference-engine-api.${var.workers_subdomain}"
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
  depends_on   = [cloudflare_record.pages]
}

# ── Workers (container — deployed via wrangler, registered here) ──────────────

resource "cloudflare_workers_script" "api" {
  account_id = var.account_id
  name       = "difference-engine-api"
  content    = "export default { async fetch(req) { return fetch(req); } }"
  module     = true
}
