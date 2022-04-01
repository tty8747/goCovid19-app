data "cloudflare_zone" "ubukubu" {
  name = var.domain
}

resource "cloudflare_record" "gitlab" {
  zone_id = data.cloudflare_zone.ubukubu.zone_id
  name    = "gitlab"
  value   = aws_eip.gitlab.public_ip
  type    = "A"
  ttl     = 1
  proxied = false
}
