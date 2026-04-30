package cnpg

import (
	kube_cnpg "danicos.dev/daniel/go-kube/pkg/cnpg"
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/homelab/pkg/root"
	pg "github.com/cloudnative-pg/api/pkg/api/v1"
	core "k8s.io/api/core/v1"
)

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
