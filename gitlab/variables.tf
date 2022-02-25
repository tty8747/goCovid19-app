variable "cloudflare_email" {}
variable "cloudflare_api_key" {}

variable "domain" {
  type    = string
  default = "ubukubu.ru"
}

variable "id_rsa_path" {
  type        = string
  default     = "~/.ssh/id_rsa.pub"
  description = "Path to public key"
}

variable "region" {
  type    = string
  default = "eu-central-1"
}

variable "cidr_block" {
  type    = string
  default = "192.168.0.0/16"
}

variable "subnet" {
  type    = string
  default = "192.168.18.0/24"
}
