# https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html
# https://docs.aws.amazon.com/eks/latest/userguide/aws-load-balancer-controller.html
resource "kubernetes_service_account" "alb_controller" {
  metadata {
    labels = {
      "app.kubernetes.io/component" = "controller"
      "app.kubernetes.io/name"      = "aws-load-balancer-controller"
    }
    name      = "aws-load-balancer-controller"
    namespace = "kube-system"
    annotations = {
      "eks.amazonaws.com/role-arn"               = "arn:aws:iam::111122223333:role/AmazonEKSLoadBalancerControllerRole"
      "eks.amazonaws.com/sts-regional-endpoints" = "true"
    }
  }
}

resource "helm_release" "alb_controller" {
  count = length(local.aws_lbc)
  name  = "alb-controller"

  # "${data.aws_caller_identity.current.id}.dkr.ecr.${data.aws_region.current.id}.amazonaws.com/amazon/aws-load-balancer-controller:2.4.0"
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
    value = kubernetes_service_account.alb_controller.metadata[0].name
  }
}

# https://docs.nginx.com/nginx-ingress-controller/installation/installation-with-helm/
# https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/ingress/annotations/
# https://aws.amazon.com/blogs/opensource/network-load-balancer-nginx-ingress-controller-eks/
resource "helm_release" "nginx_ingress" {
  name      = "nginx-ingress"
  namespace = "kube-system"

  repository = "https://helm.nginx.com/stable"
  chart      = "nginx-ingress"

  set {
    name  = "kubernetes.io/ingress.class"
    value = "alb"
  }

  set {
    name  = "alb.ingress.kubernetes.io/ip-address-type"
    value = "dualstack"
  }

  set {
    name  = "alb.ingress.kubernetes.io/scheme"
    value = "internet-facing"
  }

  set {
    name  = "alb.ingress.kubernetes.io/target-type"
    value = "ip"
  }
}

# https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html
# https://kubernetes.github.io/ingress-nginx/deploy/#aws
# resource "kubernetes_namespace" "nlb" {
#   metadata {
#     labels = {
#       Description = "Amazon_network_loadbalancer"
#     }
#     name = "nlb"
#   }
# }

# resource "kubernetes_service" "nlb" {
#   metadata {
#     name = "nlb"
#     namespace = kubernetes_namespace.nlb
#     annotations = {
#       "service.beta.kubernetes.io/aws-load-balancer-type" = "external"
#       "service.beta.kubernetes.io/aws-load-balancer-nlb-target-type" = "ip"
#       "service.beta.kubernetes.io/aws-load-balancer-scheme" = "internet-facing"
#     }
#   }
#   spec {
#     session_affinity = "None"
#     type = "LoadBalancer"
#     port {
#       port        = 80
#       target_port = 80
#       protocol = "TCP"
#     }
# 
# #   port {
# #     port        = 8080
# #     target_port = 80
# #     protocol = "TCP"
# #   }
#   }
# }

# https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/ingress_v1#tls
# resource "kubernetes_service" "app" {
#   metadata {
#     name = "ingress-service"
#   }
#   spec {
#     selector = {
#       app = "someselector"
#     }
#     port {
#       port        = 80
#       target_port = 80
#       protocol    = "TCP"
#     }
#     type = "NodePort"
#   }
# }
# 
# 
# resource "kubernetes_ingress_v1" "applb" {
#   wait_for_load_balancer = true
#   metadata {
#     name = "applb"
#     labels = {
#       app = "goCovid"
#     }
#     ## !!
#     namespace = kubernetes_namespace.nlb.id
#     annotations = {
#       # Ingress Core Settings
#       "kubernetes.io/ingress.class"      = "alb"
#       "alb.ingress.kubernetes.io/scheme" = "internet-facing"
#       # Health Check Settings
#       "alb.ingress.kubernetes.io/healthcheck-protocol"         = "HTTP"
#       "alb.ingress.kubernetes.io/healthcheck-port"             = "traffic-port"
#       "alb.ingress.kubernetes.io/healthcheck-path"             = "/usermgmt/health-status"
#       "alb.ingress.kubernetes.io/healthcheck-interval-seconds" = "15"
#       "alb.ingress.kubernetes.io/healthcheck-timeout-seconds"  = "5"
#       "alb.ingress.kubernetes.io/success-codes"                = "200"
#       "alb.ingress.kubernetes.io/healthy-threshold-count"      = "2"
#       "alb.ingress.kubernetes.io/unhealthy-threshold-count"    = "2"
#     }
#   }
#   spec {
#     rule {
#       http {
#         path {
#           path = "/*"
#           backend {
#             service {
#               name = kubernetes_service.app.metadata.0.name
#               port {
#                 number = 80
#               }
#             }
#           }
#         }
#       }
#     }
#   }
# }
