# GitOps
> A set of "best practices" where the entire code delivery process is controlled via Git,
including infrastructure and application definition as code and automation to complete updates
and rollbacks.

The key GitOps Principles
- The entire system (infrastructure and applications) is described declaratively.
- The canonical desired system state is versioned in Git.
- Changes approved are automated and applied to the system.
- Software agents ensure correctness and alert on divergence.

https://opengitops.dev/

Key points for a Kubernetes cluster:
1. The state of the cluster is always described in Git. Git holds everything for the application
and not just the source code.
2. There is no external deployment system with full access to the cluster. The cluster itself
is pulling changes and deployment information.
3. The GitOps controller is running in a constant loop and always matches the Git state with the
cluster state (reconciliation loop).

## flux vs ArgoCD
Flux is more CLI Driven
ArgoCD has a richer GUI

## Flux
Prerequisites:
- Kubernetes Cluster
- Github (giea, etc) personal access token

