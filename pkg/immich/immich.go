package immich

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	core "k8s.io/api/core/v1"
)

var Secret = struct {
	Name          string
	DBNameKey     string
	DBUsernameKey string
	DBPasswordKey string
	DBHostKey     string
	DBPortKey     string
}{
	Name:          root.Immich.Name + "-secret",
	DBNameKey:     "db_name",
	DBUsernameKey: "db_username",
	DBPasswordKey: "db_password",
	DBHostKey:     "db_host",
	DBPortKey:     "db_port",
}

var meta kube.Metadata
var Namespace = kube.Namespace(root.Immich.Name)
var srv core.Service
var uploads_pvc core.PersistentVolumeClaim

func init() {
	meta = kube.NewMetadata(root.Immich.Name, Namespace)
	uploads_pvc = meta.PVC()
	uploads_pvc.Spec.StorageClassName = new(root.TrueNASSTorageClassNFS)
	uploads_pvc.Spec.AccessModes = []core.PersistentVolumeAccessMode{core.ReadWriteMany}
}

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"namespace":   Namespace,
			"uploads-pvc": uploads_pvc,
		},
	)
	return kz.Stack("immich")
}
