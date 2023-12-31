# How to Run Examples

## Initialize and Setup Connection Variables

The Densify Terraform Provider needs to be able to connect to your instance. Port :8443 can be used just like :443 in case you have SSO (Single Sign-On) authentication enabled.

```ps1
terraform init
$env:DENSIFY_INSTANCE="https://instance.densify.com:8443"
$env:DENSIFY_USERNAME="username"
$env:DENSIFY_PASSWORD="password"
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
