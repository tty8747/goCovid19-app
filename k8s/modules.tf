module "db" {
  source         = "./modules/mariadb"
  db_engine      = "MariaDB"
  db_engineVer   = "10.5"
  db_name        = [for k, v in var.db_set : v if k == "dbname"].0
  db_user        = [for k, v in var.db_set : v if k == "dbuser"].0
  db_pass        = [for k, v in var.db_set : v if k == "dbpass"].0
  av_zone        = join("", [var.region, "b"])
  vpc_cidrs      = [module.vpc.vpc_cidr_block]
  db_subnet_list = module.vpc.private_subnets
  available_from_subnets = module.vpc.private_subnets_cidr_blocks
  vpc_id = module.vpc.vpc_id
}
