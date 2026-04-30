package cnpg

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
)

var meta kube.Metadata
var Namespace = kube.Namespace(root.CloudNativePG + "-system")
var PGClusterNamespace = kube.Namespace(root.CloudNativePG + "-cluster")

func init() {
	meta = kube.NewMetadata(root.CloudNativePG, Namespace)
}

func Stack() stack.Stack {
	s := stack.NewStack("cloud-native-pg", map[string]any{
		"namespace":              Namespace,
		"cluster-namespace":      PGClusterNamespace,
		"helm-repository-source": HelmSource(),
		"helm-release":           HelmRelease(),
		"pg-cluster":             Cluster(),
	})
	return s
}
