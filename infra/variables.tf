variable "account_id" {
  type        = string
  description = "Cloudflare account ID"
}

variable "zone_id" {
  type        = string
  description = "Cloudflare zone ID for nothinglefttolearn.com"
}

variable "domain" {
  type        = string
  description = "Domain name, e.g. nothinglefttolearn.com"
}

variable "workers_subdomain" {
  type        = string
  description = "Your workers.dev subdomain (Cloudflare dashboard → Workers & Pages → Overview)"
}
