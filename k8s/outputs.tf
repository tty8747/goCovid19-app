output "endpoint" {
  value = aws_eks_cluster.ek8s.endpoint
}

output "kubeconfig-certificate-authority-data" {
  value = aws_eks_cluster.ek8s.certificate_authority[0].data
}

output "ugrade-kube-config-command" {
  value = "aws eks update-kubeconfig --region ${data.aws_region.current.id} --name ${aws_eks_cluster.ek8s.name}"
}
