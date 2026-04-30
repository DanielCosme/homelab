package cnpg

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	v1 "k8s.io/api/core/v1"

	kube_cnpg "danicos.dev/daniel/go-kube/pkg/cnpg"
	pg "github.com/cloudnative-pg/api/pkg/api/v1"
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

func Cluster() pg.Cluster {
	meta := kube.NewMetadata(root.CloudNativePG+"-cluster", PGClusterNamespace)
	spec := pg.ClusterSpec{
		Instances: 3,
		StorageConfiguration: pg.StorageConfiguration{
			StorageClass:       new(root.KUBE_LOCAL_STORAGE_CLASS),
			Size:               "10Gi",
			ResizeInUseVolumes: new(true),
			PersistentVolumeClaimTemplate: &v1.PersistentVolumeClaimSpec{
				StorageClassName: new(root.KUBE_LOCAL_STORAGE_CLASS),
			},
		},
	}
	return kube_cnpg.NewCluster(meta, spec)
}
