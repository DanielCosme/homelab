package cnpg

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/homelab/pkg/root"

	kube_cnpg "danicos.dev/daniel/go-kube/pkg/cnpg"
	pg "github.com/cloudnative-pg/api/pkg/api/v1"
	core "k8s.io/api/core/v1"
)

func Cluster() pg.Cluster {
	meta := kube.NewMetadata(root.CloudNativePG+"-cluster", PGClusterNamespace)
	spec := pg.ClusterSpec{
		Instances: 3,
		StorageConfiguration: pg.StorageConfiguration{
			StorageClass:       new(root.KUBE_LOCAL_STORAGE_CLASS),
			Size:               "10Gi",
			ResizeInUseVolumes: new(true),
			PersistentVolumeClaimTemplate: &core.PersistentVolumeClaimSpec{
				StorageClassName: new(root.KUBE_LOCAL_STORAGE_CLASS),
			},
		},
		Managed: &pg.ManagedConfiguration{
			Roles: []pg.RoleConfiguration{},
		},
	}
	return kube_cnpg.NewCluster(meta, spec)
}

/*
 Need a role and secret
 - Username
 - Password
*/
