provider "aws" {
  alias = "ireland"
  # This resource can only be used with us-east-1 region.
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecrpublic_repository#catalog_data
  region = "us-east-1"
}

locals {
  repo_name = ["front-${var.app_name}", "back-${var.app_name}"]
}

resource "aws_ecrpublic_repository" "goCovid" {

  count           = length(local.repo_name)
  provider        = aws.ireland
  repository_name = local.repo_name[count.index]

  catalog_data {
    about_text        = "Covid tracker"
    architectures     = ["x86-64"]
    description       = "https://github.com/tty8747/goCovid19"
    logo_image_blob   = filebase64("image.png")
    operating_systems = ["Linux"]
    usage_text        = "docker pull <repo_url>/local.repo_name[count.index]"
  }
}

resource "aws_ecrpublic_repository_policy" "goCovid" {
  count           = length(local.repo_name)
  provider        = aws.ireland
  repository_name = aws_ecrpublic_repository.goCovid[count.index].repository_name

  policy = <<EOF
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Sid": "go Covid policy",
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "ecr:BatchCheckLayerAvailability",
                "ecr:PutImage",
                "ecr:InitiateLayerUpload",
                "ecr:UploadLayerPart",
                "ecr:CompleteLayerUpload",
                "ecr:DescribeRepositories",
                "ecr:GetRepositoryPolicy",
                "ecr:ListImages",
                "ecr:DeleteRepository",
                "ecr:BatchDeleteImage",
                "ecr:SetRepositoryPolicy",
                "ecr:DeleteRepositoryPolicy"
            ]
        }
    ]
}
EOF
}

data "aws_ecr_authorization_token" "token" {
  count       = length(local.repo_name)
  registry_id = aws_ecrpublic_repository.goCovid[count.index].registry_id
}

output "repo_url" {
  value = aws_ecrpublic_repository.goCovid[*].repository_uri
}
