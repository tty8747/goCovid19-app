terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.2.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.8.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "2.4.1"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "3.1.0"
    }
  }
}

provider "tls" {}

provider "aws" {
  region                   = var.region
  shared_credentials_files = ["~/.aws/credentials"]
  profile                  = "tty8747"
}

provider "null" {}

provider "kubernetes" {
  # aws eks update-kubeconfig --region eu-central-1 --name eks-myk8s-4Mqb
  # config_path    = "~/.kube/config"
  # config_context = "arn:aws:eks:${data.aws_region.current.id}:${data.aws_caller_identity.current.id}:cluster/${aws_eks_cluster.ek8s.name}"

  host                   = aws_eks_cluster.ek8s.endpoint
  cluster_ca_certificate = base64decode(aws_eks_cluster.ek8s.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.ek8s.token
}

provider "helm" {
  kubernetes {
    # config_path = "~/.kube/config"

    host                   = aws_eks_cluster.ek8s.endpoint
    cluster_ca_certificate = base64decode(aws_eks_cluster.ek8s.certificate_authority[0].data)
    token                  = data.aws_eks_cluster_auth.ek8s.token
  }
}
