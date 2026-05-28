output "stems_bucket" {
  value = cloudflare_r2_bucket.stems.name
}

output "output_bucket" {
  value = cloudflare_r2_bucket.output.name
}

output "pages_url" {
  value = "https://${var.domain}"
}

output "workers_url" {
  value = local.workers_url
}

output "r2_token" {
  value     = cloudflare_api_token.r2_worker.value
  sensitive = true
}
