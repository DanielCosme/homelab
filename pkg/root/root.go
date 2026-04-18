package root

const (
	HYDRA_CLUSTER  = "hydra"
	HYDRA_HOSTNAME = "hydra-0" // VPN Host
	GITEA_HOST     = "danicos.dev"
)

const (
	FLUX_NAMESPACE          = "flux-system"
	FLUX_APPS_HYDRA_PATH    = "./apps/" + HYDRA_CLUSTER
	FLUX_CLUSTER_HYDRA_PATH = "./clusters/" + HYDRA_CLUSTER
)
