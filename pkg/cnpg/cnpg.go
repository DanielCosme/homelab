package cnpg

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
)

var meta_system kube.Metadata
var SystemNamespace = kube.Namespace(root.CloudNativePG + "-system")

func init() {
	meta_system = kube.NewMetadata(root.CloudNativePG, SystemNamespace)
}

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta_system,
		map[string]any{
			"namespace":              SystemNamespace,
			"helm-repository-source": HelmSource(),
			"helm-release":           HelmRelease(),
		},
	)
	return kz.Stack("cloud-native-pg")
}
