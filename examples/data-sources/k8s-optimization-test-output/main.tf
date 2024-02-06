terraform {
  required_providers {
    densify = {
      source = "densify.com/provider/densify"
    }
  }
}

# credentials can be passed in as environment variables, DENSIFY_INSTANCE, DENSIFY_USERNAME, DENSIFY_PASSWORD, DENSIFY_TECH_PLATFORM, DENSIFY_ANALYSIS_NAME, DENSIFY_ENTITY_NAME

provider "densify" {
  tech_platform   = "kubernetes"
  cluster         = "<cluster-name>"
  namespace       = "<namespace>"
  controller_type = "deployment"
  pod_name        = "<pod_name>"

  # container_name  = "<container-name>"
  # continue_if_error = true
}

data "densify_container" "reco" {}

output "data_container" {
  value = data.densify_container.reco
  # value = data.densify_container.reco.containers.<container-name>.rec_cpu_req
}
