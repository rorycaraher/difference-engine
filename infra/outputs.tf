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

