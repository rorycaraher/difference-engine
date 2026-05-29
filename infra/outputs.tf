output "stems_bucket" {
  value = cloudflare_r2_bucket.stems.name
}

output "output_bucket" {
  value = cloudflare_r2_bucket.output.name
}
