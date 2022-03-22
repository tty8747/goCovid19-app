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

# https://aws.github.io/aws-eks-best-practices/cluster-autoscaling/
resource "aws_iam_role" "ek8s-AmazonEKSClusterAutoscalerRole" {
  name        = "${local.cluster_name}-AmazonEKSClusterAutoscalerRole"
  description = "Amazon EKS - Cluster autoscaler role"
  # assume_role_policy = data.aws_iam_policy_document.ek8s-AmazonEKSClusterAutoscalerPolicy.json
  assume_role_policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Effect" : "Allow",
          "Principal" : {
            "Federated" : "${aws_iam_openid_connect_provider.default.arn}"
          },
          "Action" : "sts:AssumeRoleWithWebIdentity",
          "Condition" : {
            "StringEquals" : {
              "${aws_iam_openid_connect_provider.default.url}:aud" : "sts.amazonaws.com"
            }
          }
        }
      ]
  })
}

resource "aws_iam_role_policy_attachment" "ek8s-AmazonEKSClusterAutoscalerPolicy" {
  policy_arn = aws_iam_policy.ek8s-AmazonEKSClusterAutoscalerPolicy.arn
  role       = aws_iam_role.ek8s-AmazonEKSClusterAutoscalerRole.name
}

# https://www.reddit.com/r/aws/comments/gzkzph/eksterraform_how_to_setup_aws_autoscaling_policy/
data "kubectl_file_documents" "autoscaling_yaml" {
  content = templatefile("${path.module}/cluster-autoscaler-autodiscover.yml.tftpl",
    {
      account_id                         = data.aws_caller_identity.current.account_id
      cluster_name                       = local.cluster_name
      amazon_eks_cluster_autoscaler_role = "${local.cluster_name}-AmazonEKSClusterAutoscalerRole"
    }
  )
}

# https://docs.aws.amazon.com/eks/latest/userguide/autoscaling.html
resource "kubectl_manifest" "autoscaling_yaml" {
  # Set 6 to avoid this error: The "for_each" value depends on resource attributes that cannot be determined until apply, so Terraform cannot predict how many instances will be created.
  count = 6
  # count     = length(data.kubectl_file_documents.autoscaling_yaml.documents)
  yaml_body = element(data.kubectl_file_documents.autoscaling_yaml.documents, count.index)
}
