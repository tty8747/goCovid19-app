data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

data "aws_eks_cluster_auth" "ek8s" {
  name = aws_eks_cluster.ek8s.name
}

data "tls_certificate" "ek8s" {
  url = aws_eks_cluster.ek8s.identity[0].oidc[0].issuer
}
