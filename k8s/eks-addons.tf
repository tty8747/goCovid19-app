# OpenID Connect provider
resource "aws_iam_openid_connect_provider" "default" {
  url             = aws_eks_cluster.ek8s.identity[0].oidc[0].issuer
  client_id_list  = var.openid_list
  thumbprint_list = [data.tls_certificate.ek8s.certificates[0].sha1_fingerprint]
}

# https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/
# https://github.com/kubernetes-sigs/aws-load-balancer-controller

# IAM Role for EKS Load Balancer Controller add-on
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
            curl -o iam_policy.json https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/main/docs/install/iam_policy.json
     EOT
  }
}

resource "aws_iam_policy" "ek8s-AWSLoadBalancerControllerIAMPolicy" {
  depends_on  = [null_resource.policy]
  name        = "AWSLoadBalancerControllerIAMPolicy"
  path        = "/"
  description = "AWS LoadBalancer Controller IAM Policy"

  policy = file("iam_policy.json")
}

resource "aws_iam_role" "ek8s-AmazonEKSLoadBalancerControllerRole" {
  name = "${local.cluster_name}-AmazonEKSLoadBalancerControllerRole"

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
              "${aws_iam_openid_connect_provider.default.url}:aud" : "sts.amazonaws.com",
              "${aws_iam_openid_connect_provider.default.url}:sub" : "system:serviceaccount:kube-system:aws-load-balancer-controller"
            }
          }
        }
      ]
  })
}

resource "aws_iam_role_policy_attachment" "ek8s-AWSLoadBalancerControllerIAMPolicy" {
  policy_arn = aws_iam_policy.ek8s-AWSLoadBalancerControllerIAMPolicy.arn
  role       = aws_iam_role.ek8s-AmazonEKSLoadBalancerControllerRole.name
}

resource "kubernetes_service_account" "aws_load_balancer_controller" {
  metadata {
    name      = "aws-load-balancer-controller"
    namespace = "kube-system"
    labels = {
      "app.kubernetes.io/component" = "controller"
      "app.kubernetes.io/name"      = "aws-load-balancer-controller"
    }
    annotations = {
      "eks.amazonaws.com/role-arn"               = aws_iam_role.ek8s-AmazonEKSLoadBalancerControllerRole.arn
      "eks.amazonaws.com/sts-regional-endpoints" = "true"
    }
  }
}

# Addons
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

# Amazon load balancer controller
resource "helm_release" "aws-load-balancer-controller" {
  count = length(local.aws_lbc)
  name  = "aws-load-balancer-controller"

  repository = "https://aws.github.io/eks-charts"
  chart      = local.aws_lbc[count.index]
  namespace  = "kube-system"

  set {
    name  = "clusterName"
    value = aws_eks_cluster.ek8s.name
  }

  set {
    name  = "serviceAccount.create"
    value = "false"
  }

  set {
    name  = "serviceAccount.name"
    value = kubernetes_service_account.aws_load_balancer_controller.metadata[0].name
  }
}

# NGINX Ingress as a variant
# https://docs.nginx.com/nginx-ingress-controller/installation/installation-with-helm/
# https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/ingress/annotations/
# https://aws.amazon.com/blogs/opensource/network-load-balancer-nginx-ingress-controller-eks/
# resource "helm_release" "nginx_ingress" {
#   name      = "nginx-ingress"
#   namespace = "kube-system"
#
#   repository = "https://helm.nginx.com/stable"
#   chart      = "nginx-ingress"
#
#   set {
#     name  = "kubernetes.io/ingress.class"
#     value = "alb"
#   }
#
#   set {
#     name  = "alb.ingress.kubernetes.io/ip-address-type"
#     value = "ipv4"
#   }
#
#   set {
#     name  = "alb.ingress.kubernetes.io/scheme"
#     value = "internet-facing"
#   }
#
#   set {
#     name  = "alb.ingress.kubernetes.io/target-type"
#     value = "ip"
#   }
# }
