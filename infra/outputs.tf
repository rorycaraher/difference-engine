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
  value = "https://${var.domain}/mixdown"
}

output "r2_token" {
  value     = cloudflare_api_token.r2_worker.value
  sensitive = true
}
