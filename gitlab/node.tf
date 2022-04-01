data "aws_ami" "ubuntu20" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_key_pair" "homepc" {
  key_name   = "id_rsa.pub from home pc"
  public_key = file(var.id_rsa_path)
}

resource "aws_network_interface" "gitlab" {
  subnet_id       = aws_subnet.gitlab.id
  security_groups = [aws_security_group.gitlab.id]

  tags = {
    Name = "Public network interface for gitlab instance"
  }
}

resource "aws_eip" "gitlab" {
  vpc = true

  tags = {
    Name = "gitlab"
  }
}

resource "aws_eip_association" "gitlab" {
  allocation_id        = aws_eip.gitlab.id
  network_interface_id = aws_network_interface.gitlab.id
}

resource "aws_instance" "gitlab" {
  ami = data.aws_ami.ubuntu20.id
  # t3a.medium = 0.0432 USD per Hour = 1,0368 USD per Day = 32,1408 USD per Month = 385,6896 per Year
  # t2.xlarge = 0.2144
  instance_type = "t2.xlarge"
  key_name      = aws_key_pair.homepc.id
  user_data     = data.template_file.init.rendered

  root_block_device {
    volume_size = "32"
    volume_type = "gp3"
  }

  network_interface {
    network_interface_id = aws_network_interface.gitlab.id
    device_index         = 0
  }

  tags = {
    Name = "gitlab"
  }
  depends_on = [
    aws_efs_access_point.gitlab,
  ]
}
