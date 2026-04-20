package root

import (
	core "k8s.io/api/core/v1"
)

const (
	HYDRA_CLUSTER         = "hydra"
	HYDRA_HOSTNAME        = "hydra-0" // VPN Host
	GITEA_HOST            = "danicos.dev"
	TMP_FOLDER            = "./tmp"
	SECRETS_FOLDER        = TMP_FOLDER + "/secrets"
	GO_SECRETS_FOLDER     = "./pkg/secrets"
	GO_ENC_SECRETS_FOLDER = "./pkg/enc"
)

const (
	FLUX_NAMESPACE               = "flux-system"
	FLUX_APPS_HYDRA_PATH         = "./apps/" + HYDRA_CLUSTER
	FLUX_APPS_SECRETS_HYDRA_PATH = "./apps/" + HYDRA_CLUSTER + "/secrets"
	FLUX_CLUSTER_HYDRA_PATH      = "./clusters/" + HYDRA_CLUSTER
	FLUX_DECRYPTION_PROVIDER     = "sops"
)

var (
	ContainerSecurityContext = &core.SecurityContext{
		AllowPrivilegeEscalation: new(false),
	}
)
