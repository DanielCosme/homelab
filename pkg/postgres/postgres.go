package postgres

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"

	kube_cnpg "danicos.dev/daniel/go-kube/pkg/cnpg"
	pg "github.com/cloudnative-pg/api/pkg/api/v1"
	core "k8s.io/api/core/v1"
)

var cluster_meta = kube.NewMetadata(root.CloudNativePG+"-cluster", root.PGClusterNamespace)

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		cluster_meta,
		map[string]any{
			"cluster-namespace": root.PGClusterNamespace,
			"pg-cluster":        Cluster(),
			"immich-db":         ImmichDatabase(),
		},
	)
	return kz.Stack("postgres")
}

func Cluster() pg.Cluster {
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
			Roles: []pg.RoleConfiguration{
				{
					Name:  root.Immich.Name,
					Login: true,
					PasswordSecret: &pg.LocalObjectReference{
						Name: ImmichPGSecret.Name,
					},
				},
			},
		},
	}
	return kube_cnpg.NewCluster(cluster_meta, spec)
}

var ImmichPGSecret = kube_cnpg.NewPGSecret(root.Immich.Name, root.PGClusterNamespace)

func ImmichDatabase() pg.Database {
	meta := kube.NewMetadata(root.Immich.Name+"-db", root.PGClusterNamespace)
	spec := pg.DatabaseSpec{
		Name:  root.Immich.Name,
		Owner: root.Immich.Name,
		ClusterRef: core.LocalObjectReference{
			Name: cluster_meta.Meta().Name,
		},
	}
	return kube_cnpg.NewDatabase(meta, spec)
}
