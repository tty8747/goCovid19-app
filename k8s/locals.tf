resource "random_string" "suffix" {
  length  = 4
  special = false
}

locals {
  cluster_name   = "eks-${var.k8s_name}-${random_string.suffix.result}"
  eks_node_group = "${local.cluster_name}_node_group"
  aws_lbc        = ["aws-load-balancer-controller"]
  fqdn           = join(".", [var.cname_record, var.domain])
  kubeconfig     = <<KUBECONFIG


apiVersion: v1
clusters:
- cluster:
    server: ${aws_eks_cluster.ek8s.endpoint}
    certificate-authority-data: ${aws_eks_cluster.ek8s.certificate_authority[0].data}
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: aws
  name: aws
current-context: aws
kind: Config
preferences: {}
users:
- name: aws
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      command: aws-iam-authenticator
      args:
        - "token"
        - "-i"
        - "${aws_eks_cluster.ek8s.id}"
KUBECONFIG
}
