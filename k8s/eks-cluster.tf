# IAM Role for EKS cluster
resource "aws_iam_role" "ek8s_cluster" {
  name = local.cluster_name

  assume_role_policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Effect" : "Allow",
          "Principal" : {
            "Service" : "eks.amazonaws.com"
          },
          "Action" : "sts:AssumeRole"
        }
      ]
  })
}

resource "aws_iam_role_policy_attachment" "ek8s-AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.ek8s_cluster.name
}

# Optionally, enable Security Groups for Pods
# Reference: https://docs.aws.amazon.com/eks/latest/userguide/security-groups-for-pods.html
resource "aws_iam_role_policy_attachment" "ek8s-AmazonEKSVPCResourceController" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"
  role       = aws_iam_role.ek8s_cluster.name
}

# ---
resource "aws_eks_cluster" "ek8s" {
  name                      = local.cluster_name
  role_arn                  = aws_iam_role.ek8s_cluster.arn
  version                   = "1.21"
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

# EKS logs
resource "aws_cloudwatch_log_group" "ek8s" {
  # The log group name format is /aws/eks/<cluster-name>/cluster
  # Reference: https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html
  name              = "/aws/eks/${local.cluster_name}/cluster"
  retention_in_days = 7
}
