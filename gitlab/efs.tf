resource "aws_efs_file_system" "gitlab" {
  creation_token   = "gitlab"
  performance_mode = "generalPurpose"

  tags = {
    Name = "gitlab"
  }
}

resource "aws_efs_mount_target" "gitlab" {
  file_system_id  = aws_efs_file_system.gitlab.id
  subnet_id       = aws_subnet.gitlab.id
  security_groups = [aws_security_group.efs_gitlab.id]
}

resource "aws_efs_backup_policy" "policy" {
  file_system_id = aws_efs_file_system.gitlab.id

  backup_policy {
    status = "ENABLED"
  }
}

resource "aws_efs_access_point" "gitlab" {
  file_system_id = aws_efs_file_system.gitlab.id
}
