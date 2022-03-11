output "endpoint" {
  value = aws_eks_cluster.ek8s.endpoint
}

output "kubeconfig-certificate-authority-data" {
  value = aws_eks_cluster.ek8s.certificate_authority[0].data
}

