# terraform-provider-gitsync

`terraform-provider-gitsync` is a simple Terraform provider designed to make it easy to manage `values.yaml` files directly in your Git repositories.

## Why is this useful?

When using GitOps tools such as **ArgoCD** or **FluxCD**, application deployments are typically driven by `values.yaml` files stored in Git.  
At the same time, your infrastructure is often provisioned with Terraform, which generates URLs, ARNs, IPs, and other values that are not known until after `terraform apply` completes.  
This usually forces you to manually update `values.yaml` files after the infrastructure is created.

This provider automates that process. By updating `values.yaml` as part of your Terraform run, you can fully decouple **infrastructure provisioning** from **application deployment**.

## Advantages over the traditional approach

Traditionally, you might use Terraform to create Kubernetes ConfigMaps and inject these generated values.  
With `terraform-provider-gitsync`, the values live in Git instead, which offers several benefits:

- **Visibility:** All configuration changes are tracked and reviewed through Git, without requiring direct access to the cluster.  
- **Faster recovery:** If something breaks, you can edit the `values.yaml` file directly in the repository and fix the issue immediately. Terraform can be updated later to match the new state.  
- **Safe, minimal changes:** You avoid unnecessary `terraform apply` runs that might unintentionally modify other parts of your infrastructure.

## Important note about secrets

Do **not** store sensitive values in `values.yaml`, as they are stored in plain text in Git.  
If you need to manage secrets, use a secure tool such as **SealedSecrets**, **HashiCorp Vault**, **External Secrets Operator** or a similar solution.

## Contributing

Contributions are welcome!  
There’s still plenty to do, including:

- Support for importing existing `values.yaml` files  
- Additional tests and validations  
- Improved error handling and edge case coverage
