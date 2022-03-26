data "cloudflare_zone" "ubu" {
  name = var.domain
}

resource "cloudflare_record" "cname_record" {
  zone_id = data.cloudflare_zone.ubu.zone_id
  name    = var.cname_record
  value   = aws_lb.ek8s.dns_name
  type    = "CNAME"
  ttl     = 1
  proxied = false
}
