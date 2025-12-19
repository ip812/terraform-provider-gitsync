# terraform-provider-gitsync

`terraform-provider-gitsync` is a simple Terraform provider designed to make it easy to manage `values.yaml` files directly in your Git repositories.

## Why is this useful?

When using GitOps tools such as **ArgoCD** or **FluxCD**, application deployments are typically driven by `values.yaml` files stored in Git. At the same time, your infrastructure is often provisioned with Terraform, which generates URLs, ARNs, IPs, and other values that are not known until after `terraform apply` completes. This usually forces you to manually update `values.yaml` files after the infrastructure is created. This provider automates that process. By updating `values.yaml` as part of your Terraform run, you can fully decouple **infrastructure provisioning** from **application deployment**.

## Advantages over the traditional approach

Traditionally, Terraform is used to create Kubernetes resources such as ConfigMaps and inject generated values directly into the cluster. With `terraform-provider-gitsync`, these values are stored and managed in Git instead, providing several advantages:

- **Improved visibility:** All configuration changes are versioned, reviewed, and audited in the same Git repository that contains your application manifests. There is no need to inspect Terraform state or external systems to understand the active configuration.
- **Simpler and more secure architecture:** You avoid combining multiple Terraform providers to manage Helm releases, ConfigMaps, or other Kubernetes resources directly. This removes the need to grant Terraform direct access to Kubernetes clusters, which are often located in private networks and increase security risk. Instead, Terraform only requires permission to write to the application Git repository, while a GitOps tool applies the changes to the cluster.
- **Faster incident recovery:** If something breaks, you can immediately fix the issue by editing the `values.yaml` file directly in Git. The application can recover without waiting for a Terraform run, and Terraform can later be reconciled to match the updated state.
- **Safer, minimal changes:** Reducing the number of `terraform apply` executions lowers the risk of unintentionally modifying unrelated infrastructure and keeps changes tightly scoped to application configuration.

## Important note about secrets

Do **not** store sensitive values in `values.yaml`, as they are stored in plain text in Git. If you need to manage secrets, use a secure tool such as **SealedSecrets**, **HashiCorp Vault**, **External Secrets Operator** or a similar solution.

## Contributing

Contributions are welcome! There’s still plenty to do, including:

- Additional tests and validations  
- Improved error handling and edge case coverage
