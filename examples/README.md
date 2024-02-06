# Examples

This directory contains examples that can be run/tested manually via the Terraform CLI.

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or ar testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** cloud recommendation example for the provider
* **data-sources/\*** provider examples for pulling Densify cloud & container optimization recommendations as a Terraform data-source


## Setup Connection Env Variables and Initialize Provider

The Densify Terraform Provider needs to be able to connect to your instance. Port :8443 can always be used just like :443, in case you have SSO (Single Sign-On) authentication enabled.

Windows PowerShell example:
```ps1
$env:DENSIFY_INSTANCE="https://<instance>.densify.com:8443"
$env:DENSIFY_USERNAME="<username>"
$env:DENSIFY_PASSWORD="<password>"

terraform init
```

## Plan

```ps1
terraform plan
```

## Apply

```ps1
terraform apply -auto-approve
```

## Cleanup

```ps1
terraform destroy -auto-approve
```
