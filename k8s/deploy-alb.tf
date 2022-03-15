resource "aws_lb" "ek8s" {
  name               = "alb-${local.cluster_name}"
  internal           = false
  load_balancer_type = "application"
  subnets            = module.vpc.public_subnets

  enable_cross_zone_load_balancing = true

  tags = {
    Name                       = "alb-${local.cluster_name}"
    "ingress.k8s.aws/resource" = "LoadBalancer"
    "ingress.k8s.aws/stack"    = "game-2048/ingress-2048"
    "elbv2.k8s.aws/cluster"    = "eks-myk8s-OBs0"
  }

  # lifecycle {
  #   create_before_destroy = true
  # }
}
