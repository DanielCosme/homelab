package postgres

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"

	kube_cnpg "danicos.dev/daniel/go-kube/pkg/cnpg"
	pg "github.com/cloudnative-pg/api/pkg/api/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var cluster_meta = kube.NewMetadata(root.CloudNativePG+"-cluster", root.PGClusterNamespace)
var image_catalog pg.ImageCatalog

func init() {
	image_catalog = ImageCatalog()
}

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		cluster_meta,
		map[string]any{
			"cluster-namespace": root.PGClusterNamespace,
			"pg-cluster":        Cluster(),
			"immich-db":         ImmichDatabase(),
			"image-cataglog":    image_catalog,
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
		ImageCatalogRef: &pg.ImageCatalogRef{
			TypedLocalObjectReference: core.TypedLocalObjectReference{
				APIGroup: &image_catalog.APIVersion,
				Kind:     image_catalog.Kind,
				Name:     image_catalog.Name,
			},
			Major: root.PG_VERSION_18,
		},
		PostgresConfiguration: pg.PostgresConfiguration{
			Extensions: []pg.ExtensionConfiguration{
				{
					Name: "pgvector",
				},
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

func ImageCatalog() pg.ImageCatalog {
	return pg.ImageCatalog{
		TypeMeta: meta.TypeMeta{
			Kind:       "ClusterImageCatalog",
			APIVersion: "postgresql.cnpg.io/v1",
		},
		ObjectMeta: kube.ObjectMeta(root.PG_VERSION_NAME, ""),
		Spec: pg.ImageCatalogSpec{
			Images: []pg.CatalogImage{
				{
					Major: root.PG_VERSION_18,
					Image: root.PG_VERSION_18_IMAGE,
					Extensions: []pg.ExtensionConfiguration{
						{
							Name: "pgvector",
							ImageVolumeSource: core.ImageVolumeSource{
								Reference: "ghcr.io/cloudnative-pg/pgvector:0.8.2-18-trixie",
							},
						},
					},
				},
			},
		},
	}
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
		Extensions: []pg.ExtensionSpec{
			// {
			// 	DatabaseObjectSpec: pg.DatabaseObjectSpec{
			// 		Name:   "pgvector",
			// 		Ensure: pg.EnsurePresent,
			// 	},
			// },
			// {
			// 	DatabaseObjectSpec: pg.DatabaseObjectSpec{
			// 		Name:   "cube",
			// 		Ensure: pg.EnsurePresent,
			// 	},
			// },
			// {
			// 	DatabaseObjectSpec: pg.DatabaseObjectSpec{
			// 		Name:   "earthdistance",
			// 		Ensure: pg.EnsurePresent,
			// 	},
			// },
		},
	}
	return kube_cnpg.NewDatabase(meta, spec)
}
