variable "db_engine" {}
variable "db_engineVer" {}
variable "db_name" {}
variable "db_user" {}
variable "db_pass" {}
variable "av_zone" {}
variable "vpc_cidrs" {}
variable "db_subnet_list" {}
variable "available_from_subnets" {
  description = "It is a list of subnets where database is available"
}
variable "vpc_id" {}
