data "aws_caller_identity" "current" {}
data "aws_region" "current" {}


# IAM Role for EKS cluster
resource "aws_iam_role" "ek8s_cluster" {
  name = local.cluster_name

  assume_role_policy = <<POLICY
{

  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "eks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
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


# IAM Role for EKS node group
resource "aws_iam_role" "ek8s_node_group" {
  name = local.eks_node_group

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "ek8s-AmazonEKSWorkerNodePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.ek8s_node_group.name
}

resource "aws_iam_role_policy_attachment" "ek8s-AmazonEKS_CNI_Policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.ek8s_node_group.name
}

resource "aws_iam_role_policy_attachment" "ek8s-AmazonEC2ContainerRegistryReadOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.ek8s_node_group.name
}

# IAM Role for EKS Load Balancer Controller add-on
# https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html
# https://tf-eks-workshop.workshop.aws/500_eks-terraform-workshop/575_load_balancer/tf-files.html
# https://docs.aws.amazon.com/eks/latest/userguide/aws-load-balancer-controller.html
resource "null_resource" "policy" {
  triggers = {
    always_run = timestamp()
  }
  provisioner "local-exec" {
    on_failure  = fail
    when        = create
    interpreter = ["/bin/bash", "-c"]
    command     = <<EOT
            curl -o iam-policy.json https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/main/docs/install/iam_policy.json
     EOT
  }
}

resource "aws_iam_policy" "load-balancer-policy" {
  depends_on  = [null_resource.policy]
  name        = "AWSLoadBalancerControllerIAMPolicy"
  path        = "/"
  description = "AWS LoadBalancer Controller IAM Policy"

  policy = file("iam-policy.json")
}

# IAM Role for OpenID Connect provider
resource "aws_iam_openid_connect_provider" "default" {
  url             = aws_eks_cluster.ek8s.identity[0].oidc[0].issuer
  client_id_list  = var.openid_list
  thumbprint_list = []
}

resource "aws_iam_role" "AmazonEKSLoadBalancerControllerRole" {
  name = "${local.cluster_name}-AmazonEKSLoadBalancerControllerRole"

  assume_role_policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Effect" : "Allow",
          "Principal" : {
            "Federated" : "arn:aws:iam::${data.aws_caller_identity.current.account_id}:oidc-provider/${aws_iam_openid_connect_provider.default.url}"
          },
          "Action" : "sts:AssumeRoleWithWebIdentity",
          "Condition" : {
            "StringEquals" : {
              "${aws_iam_openid_connect_provider.default.url}:aud" : "sts.amazonaws.com",
              "${aws_iam_openid_connect_provider.default.url}:sub" : "system:serviceaccount:kube-system:aws-load-balancer-controller"
            }
          }
        }
      ]
  })
}


resource "aws_iam_role_policy_attachment" "ek8s-AWSLoadBalancerControllerIAMPolicy" {
  policy_arn = "arn:aws:iam::${data.aws_caller_identity.current.id}:policy/AWSLoadBalancerControllerIAMPolicy"
  role       = aws_iam_role.AmazonEKSLoadBalancerControllerRole.name
}
