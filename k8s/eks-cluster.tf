data "aws_eks_cluster_auth" "ek8s" {
  name = local.cluster_name
}

resource "aws_eks_cluster" "ek8s" {
  name     = local.cluster_name
  role_arn = aws_iam_role.ek8s_cluster.arn

  enabled_cluster_log_types = ["api"]

  vpc_config {
    subnet_ids = module.vpc.private_subnets
  }

  # Ensure that IAM Role permissions are created before and deleted after EKS Cluster handling.
  # Otherwise, EKS will not be able to properly delete EKS managed EC2 infrastructure such as Security Groups.
  depends_on = [
    aws_iam_role_policy_attachment.ek8s-AmazonEKSClusterPolicy,
    aws_iam_role_policy_attachment.ek8s-AmazonEKSVPCResourceController,
  ]
}

resource "aws_eks_node_group" "ek8s" {
  cluster_name    = aws_eks_cluster.ek8s.name
  node_group_name = local.eks_node_group
  node_role_arn   = aws_iam_role.ek8s_node_group.arn
  subnet_ids      = module.vpc.private_subnets

  # ubuntu ami types -> https://cloud-images.ubuntu.com/aws-eks/
  ami_type = "BOTTLEROCKET_x86_64"
  # t2.micro - free tier
  # instance_types = ["t2.micro"]

  scaling_config {
    desired_size = 4
    max_size     = 10
    min_size     = 2
  }

  update_config {
    max_unavailable = 1
  }

  # Ensure that IAM Role permissions are created before and deleted after EKS Node Group handling.
  # Otherwise, EKS will not be able to properly delete EC2 Instances and Elastic Network Interfaces.
  depends_on = [
    aws_iam_role_policy_attachment.ek8s-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.ek8s-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.ek8s-AmazonEC2ContainerRegistryReadOnly,
  ]

  tags = {
    "k8s.io/cluster-autoscaler/${local.cluster_name}" = "owned"
    "k8s.io/cluster-autoscaler/enabled"               = "TRUE"
  }
}

resource "aws_eks_addon" "vpc_cni" {
  cluster_name      = aws_eks_cluster.ek8s.name
  addon_name        = "vpc-cni"
  resolve_conflicts = "OVERWRITE"

  depends_on = [
    aws_iam_role_policy_attachment.ek8s-AmazonEKS_CNI_Policy,
  ]
}

resource "aws_eks_addon" "coredns" {
  cluster_name      = aws_eks_cluster.ek8s.name
  addon_name        = "coredns"
  resolve_conflicts = "OVERWRITE"
}

resource "aws_eks_addon" "kube-proxy" {
  cluster_name      = aws_eks_cluster.ek8s.name
  addon_name        = "kube-proxy"
  resolve_conflicts = "OVERWRITE"
}

# EKS logs
resource "aws_cloudwatch_log_group" "ek8s" {
  # The log group name format is /aws/eks/<cluster-name>/cluster
  # Reference: https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html
  name              = "/aws/eks/${local.cluster_name}/cluster"
  retention_in_days = 7
}
