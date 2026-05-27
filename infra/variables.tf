variable "account_id" {
  type        = string
  description = "Cloudflare account ID"
}

variable "zone_id" {
  type        = string
  description = "Cloudflare zone ID for the domain"
}

variable "domain" {
  type        = string
  description = "Domain name, e.g. difference-engine.com"
}
