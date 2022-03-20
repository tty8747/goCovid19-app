# resource "helm_release" "kuber_state_metrics" {
#   name      = "metrics-server"
#   namespace = "kube-system"
# 
#   repository = "https://prometheus-community.github.io/helm-charts"
#   chart      = "kube-state-metrics"
# }

data "kubectl_file_documents" "metrics-server" {
  content = file("metrics-server.yml")
}

resource "kubectl_manifest" "metrics-server" {
  for_each  = data.kubectl_file_documents.metrics-server.manifests
  yaml_body = each.value
}
