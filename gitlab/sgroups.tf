resource "aws_security_group" "gitlab" {
  name        = "gitlab"
  description = "Allow all"
  vpc_id      = aws_vpc.gitlab.id

  tags = {
    Name = "gitlab"
  }
}

resource "aws_security_group_rule" "ingress_all" {
  type              = "ingress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.gitlab.id
}

resource "aws_security_group_rule" "egress_all" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.gitlab.id
}

resource "aws_security_group" "efs_gitlab" {
  name   = "efs_gitlab"
  vpc_id = aws_vpc.gitlab.id

  ingress {
    description = "efs from VPC"
    from_port   = 2049
    to_port     = 2049
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.gitlab.cidr_block]
  }

  egress {
    from_port   = 2049
    to_port     = 2049
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.gitlab.cidr_block]
  }

  tags = {
    Name = "efs gitlab"
  }
}
