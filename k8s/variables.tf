variable "id_rsa_path" {
  type        = string
  default     = "~/.ssh/id_rsa.pub"
  description = "Path to public key"
}

variable "region" {
  type    = string
  default = "eu-central-1"
}

variable "cidr" {
  type        = string
  default     = "10.0.0.0/16"
  description = "CIDR for k8s vpc"
}

variable "priv_subnets" {
  type    = list(string)
  default = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "pub_subnets" {
  type    = list(string)
  default = ["10.0.11.0/24", "10.0.12.0/24", "10.0.13.0/24"]
}

variable "k8s_name" {
  type    = string
  default = "myk8s"
}

variable "openid_list" {
  type    = list(string)
  default = ["sts.amazonaws.com"]
}

variable "domain" {
  type    = string
  default = "ubukubu.ru"
}

variable "cname_record" {
  type    = string
  default = "app"
}

variable "cloudflare_email" {
  type = string
}

variable "cloudflare_api_key" {
  type = string
}

variable "db_set" {
  type = map(string)
  default = {
    "dbname" = "test"
    "dbuser" = "someuser"
    "dbpass" = "somepass"
  }
}

variable "environments" {
  type = list(string)
  default = ["test", "dev", "release"]
}
