resource "aws_vpc" "gitlab" {
  cidr_block           = var.cidr_block
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "gitlab"
  }
}

resource "aws_subnet" "gitlab" {
  vpc_id                  = aws_vpc.gitlab.id
  cidr_block              = var.subnet
  map_public_ip_on_launch = "true"
  availability_zone       = var.subnet_zone

  tags = {
    Name = "gitlab"
  }
}

resource "aws_internet_gateway" "gitlab" {
  vpc_id = aws_vpc.gitlab.id

  tags = {
    Name = "gitlab"
  }
}

resource "aws_default_route_table" "public" {
  default_route_table_id = aws_vpc.gitlab.default_route_table_id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gitlab.id
  }

  tags = {
    Name = "gitlab"
  }
}

resource "aws_route53_zone" "private" {
  name = "gitlab.company.lan"

  vpc {
    vpc_id = aws_vpc.gitlab.id
  }
}
