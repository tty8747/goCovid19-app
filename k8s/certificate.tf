resource "aws_acm_certificate" "app" {
  domain_name       = local.fqdn
  validation_method = "DNS"

  subject_alternative_names = [
    "www.${local.fqdn}",
  ]

}

# resource "cloudflare_record" "foobar" {
#   zone_id = data.cloudflare_zone.ubu.zone_id
#   name    = "terraform"
#   value   = "192.168.0.11"
#   type    = "A"
#   ttl     = 3600
# }
