package truenas

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
)

var meta kube.Metadata
var Namespace = kube.Namespace(root.TrueNAS_CSI)

func init() {
	meta = kube.NewMetadata(root.TrueNAS_CSI, Namespace)
}

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"namespace": Namespace,
		},
	)
	return kz.Stack(root.TrueNAS_CSI)
}
