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
  cluster         = "k8master"
  namespace       = "qa-llc"
  controller_type = "deployment"
  pod_name        = "webserver-deployment"

  # container_name  = "den-web"
  # skip_errors = true
}

data "densify_container" "reco" {}

output "data_container" {
  value = data.densify_container.reco
  # value = data.densify_container2.reco.containers.den-web.rec_cpu_req
}
