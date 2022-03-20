output "endpoint" {
  value = aws_eks_cluster.ek8s.endpoint
}

output "kubeconfig-certificate-authority-data" {
  value = aws_eks_cluster.ek8s.certificate_authority[0].data
}

output "upgrade-kube-config-command" {
  value = "aws eks update-kubeconfig --region ${data.aws_region.current.id} --name ${aws_eks_cluster.ek8s.name}"
}

output "db_endpoint" {
  value = module.db.db_endpoint
}

output "alb_name" {
  value = aws_lb.ek8s.name
}
