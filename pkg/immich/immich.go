package immich

import (
	"fmt"

	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/postgres"
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
	Name:          postgres.ImmichPGSecret.Name,
	DBNameKey:     "db_name",
	DBUsernameKey: "db_username",
	DBPasswordKey: "db_password",
	DBHostKey:     "db_host",
	DBPortKey:     "db_port",
}

const RedisPort = 6379
const MachineLearningPort = 3003

var meta kube.Metadata
var Namespace = kube.Namespace(root.Immich.Name)
var srv core.Service
var uploads_pvc core.PersistentVolumeClaim
var redis_meta kube.Metadata
var redis_srv core.Service
var machine_learning_meta kube.Metadata
var machine_learning_srv core.Service
var machine_learning_pvc core.PersistentVolumeClaim

func init() {
	meta = kube.NewMetadata(root.Immich.Name, Namespace)
	srv = meta.Service(root.Immich.Port)
	srv.Spec.Type = core.ServiceTypeNodePort
	srv.Spec.Ports[0].NodePort = int32(root.Immich.Public.NodePort)

	uploads_pvc = meta.PVC()
	uploads_pvc.Spec.StorageClassName = new(root.TrueNASSTorageClassNFS)
	uploads_pvc.Spec.AccessModes = []core.PersistentVolumeAccessMode{core.ReadWriteMany}
	uploads_pvc.Spec.Resources = core.VolumeResourceRequirements{
		Requests: kube.StorageRequest(resource.MustParse("200Gi")),
	}

	redis_meta = kube.NewMetadata("redis", Namespace)
	redis_srv = redis_meta.Service(RedisPort)

	machine_learning_meta = kube.NewMetadata(root.Immich.Name+"-machine-learning", Namespace)
	machine_learning_srv = machine_learning_meta.Service(MachineLearningPort)
	machine_learning_pvc = machine_learning_meta.PVC()
	machine_learning_pvc.Spec.StorageClassName = new(root.LONGHORN_STORAGE_CLASS)
}

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"namespace":            Namespace,
			"uploads-pvc":          uploads_pvc,
			"redis-srv":            redis_srv,
			"redis-deployment":     Redis(),
			"immich-deployment":    Deployment(),
			"immich-srv":           srv,
			"machine-learning-pvc": machine_learning_pvc,
			"machine-learning-srv": machine_learning_srv,
			"machine-learning":     MachineLearning(),
		},
	)
	return kz.Stack("immich")
}

func MachineLearning() apps.Deployment {
	cacheVol := kube.NewVolumeFrom(kube.VolumeSourcePVC, "cache", machine_learning_pvc.Name)
	podSpec := core.PodSpec{
		Containers: []core.Container{
			{
				Name:  machine_learning_meta.Meta().Name,
				Image: "ghcr.io/immich-app/immich-machine-learning:v2.7.5",
				VolumeMounts: []core.VolumeMount{{
					Name:      cacheVol.Name,
					MountPath: "/cache",
				}},
			},
		},
		Volumes: []core.Volume{cacheVol},
	}
	return kube.NewDeployment(machine_learning_meta, podSpec)
}

func Deployment() apps.Deployment {
	uploadVol := kube.NewVolumeFrom(kube.VolumeSourcePVC, "upload", uploads_pvc.Name)
	localtimeVol := core.Volume{
		Name: "localtime",
		VolumeSource: core.VolumeSource{
			HostPath: &core.HostPathVolumeSource{
				Path: "/etc/localtime",
				Type: new(core.HostPathFile),
			},
		},
	}
	envMapping := map[string]string{
		"REDIS_PORT":  fmt.Sprintf("%d", RedisPort),
		"IMMICH_PORT": fmt.Sprintf("%d", root.Immich.Port),
	}
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
				VolumeMounts: []core.VolumeMount{
					{
						Name:      uploadVol.Name,
						MountPath: "/data",
					},
					{
						Name:      localtimeVol.Name,
						MountPath: "/etc/localtime",
						ReadOnly:  true,
					},
				},
			},
		},
		Volumes: []core.Volume{
			uploadVol,
			localtimeVol,
		},
	}
	return kube.NewDeployment(meta, podSpec)
}

func Redis() apps.Deployment {
	podSpec := core.PodSpec{
		Containers: []core.Container{
			{
				Name:  "redis",
				Image: "docker.io/valkey/valkey:9@sha256:3b55fbaa0cd93cf0d9d961f405e4dfcc70efe325e2d84da207a0a8e6d8fde4f9",
			},
		},
	}
	return kube.NewDeployment(redis_meta, podSpec)
}
