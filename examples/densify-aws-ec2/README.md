# How to run

```ps1
$env:DENSIFY_INSTANCE="https://instance.densify.com:443"
$env:DENSIFY_USERNAME="username"
$env:DENSIFY_PASSWORD="password"
terraform plan
```

Then
```ps1
terraform apply -auto-approve
```

## Troubleshooting
To set up logging, set the TF_LOG environment variable to "INFO/DEBUG/TRACE"
```ps1
$env:TF_LOG="DEBUG"
```
