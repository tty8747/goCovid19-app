# IAM Role for Autoscalling
resource "aws_iam_policy" "ek8s-AmazonEKSClusterAutoscalerPolicy" {
  name = "${local.cluster_name}-AmazonEKSClusterAutoscalerPolicy"

  policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : [
            "autoscaling:DescribeAutoScalingGroups",
            "autoscaling:DescribeAutoScalingInstances",
            "autoscaling:DescribeLaunchConfigurations",
            "autoscaling:DescribeTags",
            "autoscaling:SetDesiredCapacity",
            "autoscaling:TerminateInstanceInAutoScalingGroup",
            "ec2:DescribeLaunchTemplateVersions"
          ],
          "Resource" : "*",
          "Effect" : "Allow"
        }
      ]
  })
}

data "aws_iam_policy_document" "ek8s-AmazonEKSClusterAutoscalerPolicy" {
  statement {
    sid = "2"

    actions = ["sts:AssumeRoleWithWebIdentity"]
    effect  = "Allow"

    condition {
      test     = "StringEquals"
      variable = "${aws_iam_openid_connect_provider.default.url}:sub"
      values   = ["system:serviceaccount:kube-system:aws-node"]
    }

    principals {
      identifiers = [aws_iam_openid_connect_provider.default.arn]
      type        = "Federated"
    }
  }
}

resource "aws_iam_role" "ek8s-AmazonEKSClusterAutoscalerRole" {
  name               = "${local.cluster_name}-AmazonEKSClusterAutoscalerRole"
  assume_role_policy = data.aws_iam_policy_document.ek8s-AmazonEKSClusterAutoscalerPolicy.json
}

resource "aws_iam_role_policy_attachment" "ek8s-AmazonEKSClusterAutoscalerPolicy" {
  policy_arn = aws_iam_policy.ek8s-AmazonEKSClusterAutoscalerPolicy.arn
  role       = aws_iam_role.ek8s-AmazonEKSClusterAutoscalerRole.name
}
