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
	KUBE_LOCAL_STORAGE_CLASS = "local-path"
)

const (
	FLUX_NAMESPACE                = "flux-system"
	FLUX_APPS_HYDRA_PATH          = "./apps/" + HYDRA_CLUSTER
	FLUX_APPS_SECRETS_HYDRA_PATH  = "./apps/" + HYDRA_CLUSTER + "/secrets"
	FLUX_CLUSTER_HYDRA_PATH       = "./clusters/" + HYDRA_CLUSTER
	FLUX_INFRA_HYDRA_PATH         = "./infrastructure/" + HYDRA_CLUSTER
	FLUX_DECRYPTION_PROVIDER      = "sops"
	FLUX_HELM_MONITORING_INTERVAL = 24 // in hours
)

const (
	HELM_PROMETHEUS_URL                = "https://prometheus-community.github.io/helm-charts"
	HELM_PROMETHEUS_CHART              = "kube-prometheus-stack"
	HELM_PROMETHEUS_CHART_VERSION      = "66.x"
	HELM_LONGHORN_URL                  = "https://charts.longhorn.io"
	HELM_LONGHORN_CHART                = "longhorn"
	HELM_LONGHORN_CHART_VERSION        = "1.11.1"
	HELM_CLOUD_NATIVE_PG_URL           = "https://cloudnative-pg.github.io/charts"
	HELM_CLOUD_NATIVE_PG_CHART         = "cloudnative-pg"
	HELM_CLOUD_NATIVE_PG_CHART_VERSION = "0.28.0"
	PG_VERSION_NAME                    = "pg-minimal-trixie"
	PG_VERSION_18                      = 18
	PG_VERSION_18_IMAGE                = "ghcr.io/cloudnative-pg/postgresql:18.3-minimal-trixie"
)

var (
	ContainerSecurityContext = &core.SecurityContext{
		AllowPrivilegeEscalation: new(false),
	}
	HYDRA_WORKERS = []string{
		"hydra-1",
		"hydra-2",
	}
)
