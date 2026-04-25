## Homelab

## GitOps
A set of "best practices" where the entire code delivery process is controlled via Git, including infrastructure and application definition as code and automation to complete updates and rollbacks.

The key GitOps Principles:

- The entire system (infrastructure and applications) is described declaratively.
- The canonical desired system state is versioned in Git.
- Changes approved are automated and applied to the system.
- Software agents ensure correctness and alert on divergence.
https://opengitops.dev/

Key points for a Kubernetes cluster:

The state of the cluster is always described in Git. Git holds everything for the application and not just the source code.
There is no external deployment system with full access to the cluster. The cluster itself is pulling changes and deployment information.
The GitOps controller is running in a constant loop and always matches the Git state with the cluster state (reconciliation loop).
flux vs ArgoCD
Flux is more CLI Driven ArgoCD has a richer GUI

## Dev dependencies
- AGE
- SOPS
- Go Toolchain

## Secrets
Some secrets are written a Go files, they must be decrypted first.
Once decrypted they show up at `./pkg/secrets/..`. Note that they are referenced in `./cmd/secrets/main.go`.

## Flux
Prerequisites:
- Kubernetes Cluster
-Github (giea, etc) personal access token

### Managing secrets with SOPS
First step to manage secrets in Flux with SOPS, is to add the secret-key to Kubernetes and make it accessible to the flux system.

To decrypt secrets use:
```sh
sops --decrypt <file_path>
```
There must be an environment variable named `SOPS_AGE_KEY_FILE` that contains the path of the AGE key text file.

## Pending
### Storage
Object Storage from NAS
Block Storage from NAS

Distributed Storage from:
	- RookFS
	- Longhorn
