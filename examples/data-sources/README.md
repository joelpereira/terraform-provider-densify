# How to Run Examples

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
