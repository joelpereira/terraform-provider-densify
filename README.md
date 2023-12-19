<img src="https://www.densify.com/wp-content/uploads/densify.png" width="300">

# Densify Terraform Provider
This provider interfaces between Densify machine learning analytics and Terraform templates.
It enables two operations:
- automated optimization of instance families/sizes, and container resource requests/limits (making them “self-optimizing”)
- auto-tagging of cloud instances and containers based on Densify’s optimization analysis (making them “self-aware”)

This integration is based on the Densify SaaS engine API which contains operational intelligence, analysis findings, and optimization recommendations for each cloud instance or container in scope.
The result is next-generation resource optimization with the elimination of hard-coded resource specifications.

- [Requirements](#requirements)
- [Usage](#usage)
- [Documentation](#docs)
- [Examples](#examples)
- [Inputs](#inputs)
- [Outputs](#outputs)
- [License](#license)

## Requirements
- Densify service account, which is provided with a Densify subscription or free trial (www.densify.com/product/trial)

## Usage

```hcl
terraform {
  required_providers {
    densify = {
      source = "densify.com/provider/densify"
    }
  }
}

provider "densify" {
  # credentials and other parameters can be passed in to the Densify Provider as environment variables

  tech_platform = "aws"
  account_number = "9876543210"
  system_name = "system-name-321"
}
```

Then we can simply utilize the optimized instance that Densify recommended:
```hcl
...
instance_type = densify.data.approved_instance
...
```

### Data Sources
There are two data sources available within the Densify Provider:
| Name | Description | Call |
|------|-------------|:-------:|
| Cloud Recommendation | This returns one cloud (AWS/Azure/GCP) recommendation from Densify | _cloud |
| Container Recommendation | This returns one container (Kubernetes) recommendation from Densify | _container |

## Documentation

You can find the generated documentation in the [docs folder](docs/).

## Examples 
* [Cloud Optimization Test Output](examples/data-sources/cloud-optimization-test-output)
* [AWS EC2](examples/data-sources/aws-ec2)
* [Azure VM](examples/data-sources/azure-vm)
* [Kubernetes Deployment](examples/data-sources/k8s-deployment)
* [Kubernetes Optimization Test Output](examples/data-sources/k8s-optimization-test-output)

## Inputs

### Densify Cloud Recommendation
Inputs for "_cloud" provider call are:

| Name | Description | Type | Environment Variable | Required |
|------|-------------|:----:|:--------------------:|:--------:|
| densify_instance | Your Densify SaaS instance URL to pull recommendations | String | DENSIFY_INSTANCE | Yes |
| username | Densify service account user name (you can request one by contacting your Account Manager or support@densify.com) | String | DENSIFY_USERNAME | Yes |
| password | Densify service account password  | String | DENSIFY_PASSWORD | Yes |
| tech_platform | The technology platform or CSP (cloud service provider) being used. Select one of the following options: aws, azure, gcp, kubernetes. | String | DENSIFY_TECH_PLATFORM | Yes |
| account_number | description | String | DENSIFY_ACCOUNT_NUMBER | Yes for cloud instances |
| system_name | description | String | DENSIFY_SYSTEM_NAME | Yes for cloud instances |
| fallback | The fallback/default instance type | String | DENSIFY_FALLBACK | No |
| skip_errors | The fallback/default instance type | String | DENSIFY_SKIP_ERRORS | No |


### Densify Container Recommendation
Inputs for "_container" provider call are:

| Name | Description | Type | Environment Variable | Required |
|------|-------------|:----:|:--------------------:|:--------:|
| densify_instance | Your Densify SaaS instance URL to pull recommendations | String | DENSIFY_INSTANCE | Yes |
| username | Densify service account user name (you can request one by contacting your Account Manager or support@densify.com) | String | DENSIFY_USERNAME | Yes |
| password | Densify service account password  | String | DENSIFY_PASSWORD | Yes |
| tech_platform | The technology platform or CSP (cloud service provider) being used. Select one of the following options: aws, azure, gcp, kubernetes. | String | DENSIFY_TECH_PLATFORM | Yes |
| cluster | description | String | DENSIFY_CLUSTER | Yes for container recommendations |
| namespace  | description | String | DENSIFY_NAMESPACE | Yes for container recommendations |
| controller_type | description | String | DENSIFY_CONTROLLER_TYPE | Yes for container recommendations |
| pod_name | description | String | DENSIFY_POD_NAME | Yes for container recommendations |
| container_name | description | String | DENSIFY_CONTAINER_NAME | Yes for container recommendations |
| fallback | The fallback/default instance type | String | DENSIFY_FALLBACK | No |
| skip_errors | The fallback/default instance type | String | DENSIFY_SKIP_ERRORS | No |


## Outputs

### Densify Cloud Recommendation
Outputs for "_cloud" provider call are:

| Name | Type | Description |
|------|------|-------------|
| entity_id | String | Unique identifier for cloud resource. |
| name | String | System name for the compute resource. |
| current_type | String | Current instance type. |
| recommended_type | String | Recommended instance type generated by Densify. |
| approved_type | String | The approved instance type. This starts with the fallback instance or the current instance type, and may only be replaced by the recommended instance if 'Approval_Type' is set. |
| optimization_type | String | Type of optimization. Ex. Downsize, Upsize, Terminate, etc. |
| account_id | String | Account reference identifier. |
| approval_type | String | Approval type. If ITSM integration has been enabled, this field will identify whether the recommendation has been reviewed & approved. |
| savings_estimate | Float64 | Estimated monthly savings by applying the optimization recommendation. |
| effort_estimate | String | Estimated effort required by applying optimization recommendation. Ex. none, low, med, high. |

### Densify Container Recommendation
Outputs for "_container" provider call are:

| Name | Description |
|------|-------------|
| entity_id | String | Unique identifier for cloud resource. |
| name | String | Container manifest name. |
| optimization_type | String | Type of optimization. Ex. Downsize, Upsize, Terminate, etc. |
| account_id | String | Account reference identifier. |
| approval_type | String | Approval type. If ITSM integration has been enabled, this field will identify whether the recommendation has been reviewed & approved. |
| cluster | String | desc |
| namespace | String | desc |
| controller_type | String | desc |
| pod_name | String | desc |
| container_name | String | desc |
| current_cpu_req | String | The current CPU Request for resources (in millicores or m). |
| current_cpu_limit | String | The current CPU Limit for resources (in millicores or m). |
| current_mem_req | String | The current Memory Request for resources (in mebibytes or Mi). |
| current_mem_limit | String | The current Memory Limit for resources (in mebibytes or Mi). |
| recommended_cpu_req | String | The recommended CPU Request for resources (in millicores or m). |
| recommended_cpu_limit | String | The recommended CPU Limit for resources (in millicores or m). |
| recommended_mem_req | String | The recommended Memory Request for resources (in mebibytes or Mi). |
| recommended_mem_limit | String | The recommended Memory Limit for resources (in mebibytes or Mi). |


## License

Apache 2 Licensed. See LICENSE for full details.


# How to Build the Densify Terraform Provider

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install .
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

