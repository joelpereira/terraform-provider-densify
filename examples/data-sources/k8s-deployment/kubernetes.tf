provider "kubernetes" {
  host                   = var.k8s_server
  cluster_ca_certificate = base64decode(var.k8s_certificate_authority_data)
  client_certificate     = base64decode(var.k8s_client_certificate_data)
  client_key             = base64decode(var.k8s_client_key_data)
}
