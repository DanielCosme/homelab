package immich

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
	srv = meta.Service(root.Immich.Port)
	uploads_pvc = meta.PVC()
	uploads_pvc.Spec.StorageClassName = new(root.TrueNASSTorageClassNFS)
	uploads_pvc.Spec.AccessModes = []core.PersistentVolumeAccessMode{core.ReadWriteMany}
	req := kube.StorageRequest(resource.MustParse("200Gi"))
	uploads_pvc.Spec.Resources = core.VolumeResourceRequirements{
		Requests: req,
	}
}

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"namespace":   Namespace,
			"uploads-pvc": uploads_pvc,
			// "redis":       Redis(),
			"deployment": Deployment(),
		},
	)
	return kz.Stack("immich")
}

func Deployment() apps.Deployment {
	uploadVol := kube.NewVolumeFrom(kube.VolumeSourcePVC, "upload", uploads_pvc.Name)
	envMapping := map[string]string{}
	secretMapping := map[string]string{
		"DB_HOSTNAME":      Secret.DBHostKey,
		"DB_PORT":          Secret.DBPortKey,
		"DB_USERNAME":      Secret.DBUsernameKey,
		"DB_PASSWORD":      Secret.DBPasswordKey,
		"DB_DATABASE_NAME": Secret.DBUsernameKey,
	}
	podSpec := core.PodSpec{
		Containers: []core.Container{
			{
				Name:  root.Immich.Name + "-server",
				Image: root.Immich.Image,
				Env:   kube.NewEnvVarWithSecret(envMapping, secretMapping, Secret.Name),
				Ports: []core.ContainerPort{{ContainerPort: root.Immich.Port}},
			},
		},
		Volumes: []core.Volume{
			// /etc/localtime
			uploadVol,
		},
	}
	return kube.NewDeployment(meta, podSpec)
}

func Redis() apps.Deployment {
	meta := kube.NewMetadata(root.Immich.Name+"-redis", Namespace)
	podSpec := core.PodSpec{}
	return kube.NewDeployment(meta, podSpec)
}
