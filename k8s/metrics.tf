resource "helm_release" "kuber_state_metrics" {
  name      = "kube-state-metrics"
  namespace = "kube-system"

  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-state-metrics"
}
