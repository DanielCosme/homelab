package immich

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	core "k8s.io/api/core/v1"
)

var Secret = struct {
	Name           string
	DB_UserKey     string
	DB_PasswordKey string
}{
	Name:           root.Immich.Name,
	DB_UserKey:     "db_username",
	DB_PasswordKey: "db_password",
}

var meta kube.Metadata
var Namespace = kube.Namespace(root.Immich.Name)
var srv core.Service
var pvc core.PersistentVolumeClaim

func init() {
	meta = kube.NewMetadata(root.Immich.Name, Namespace)
}

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"namespace": Namespace,
		},
	)
	return kz.Stack("immich")
}
