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
