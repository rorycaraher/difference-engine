# R2 buckets
resource "cloudflare_r2_bucket" "stems" {
  account_id = var.account_id
  name       = "difference-engine-stems"
  location   = "EEUR"
}

resource "cloudflare_r2_bucket" "output" {
  account_id = var.account_id
  name       = "difference-engine-output"
  location   = "EEUR"
}
