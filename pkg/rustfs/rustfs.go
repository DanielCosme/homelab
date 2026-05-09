package rustfs

import (
	"fmt"

	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

var meta kube.Metadata
var Namespace = kube.Namespace(root.RUSTFS)
var srv core.Service
var pvc core.PersistentVolumeClaim

var Secret = struct {
	Name      string
	AccessKey string
	SecretKey string
}{
	Name:      root.RUSTFS + "-secret",
	AccessKey: "access-key",
	SecretKey: "secret-key",
}

func init() {
	meta = kube.NewMetadata(root.RUSTFS, Namespace)
	srv = meta.ServiceFrom([]kube.ServicePort{
		{
			Name:     "console",
			Port:     root.RUSTFS_CONSOLE_PORT,
			NodePort: root.RUSTFS_CONSOLE_NODE_PORT,
		},
		{
			Name:     "server",
			Port:     root.RUSTFS_SERVER_PORT,
			NodePort: root.RUSTFS_CONSOLE_NODE_PORT,
		},
	}...)
	srv.Spec.Type = core.ServiceTypeNodePort

	pvc = meta.PVCFrom(kube.StorageRequests20Gi)
	pvc.Spec.StorageClassName = new(root.TrueNASSTorageClass_iSCSI)
	// TODO: restrict privledge escalation in pod
	//
	// s3.danicos.dev (API server)
	// rustfs.danicos.dev (Console)
	// Volume EmptyDir for logs
}

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"namespace": Namespace,
			"service":   srv,
			// "pvc":        pvc,
			// "deployment": Deployment(),
		},
	)
	return kz.Stack("rustfs")
}

func Deployment() apps.Deployment {
	mountPath := "/data/rustfs0"
	envMap := map[string]string{
		"RUSTFS_VOLUMES":                      mountPath,
		"RUSTFS_ADDRESS":                      fmt.Sprintf("0.0.0.0:%d", root.RUSTFS_SERVER_PORT),
		"RUSTFS_CONSOLE_ADDRESS":              fmt.Sprintf("0.0.0.0:%d", root.RUSTFS_CONSOLE_PORT),
		"RUSTFS_CONSOLE_ENABLE":               "true",
		"RUSTFS_CONSOLE_CORS_ALLOWED_ORIGINS": "*",
		"RUSTFS_OBS_LOGGER_LEVEL":             "info",
	}
	envMapSecret := map[string]string{
		"RUSTFS_ACCESS_KEY": Secret.AccessKey, // Username
		"RUSTFS_SECRET_KEY": Secret.SecretKey, // Password
	}
	dataVol := kube.NewVolumeFrom(kube.VolumeSourcePVC, "data0", pvc.Name)
	logsVol := core.Volume{
		Name: "logs",
		VolumeSource: core.VolumeSource{
			EmptyDir: &core.EmptyDirVolumeSource{},
		},
	}
	podSpec := core.PodSpec{
		Containers: []core.Container{{
			Name:  root.RUSTFS,
			Image: root.RUSTFS_IMAGE,
			Env:   kube.NewEnvVarWithSecret(envMap, envMapSecret, Secret.Name),
			Ports: []core.ContainerPort{
				{ContainerPort: root.RUSTFS_SERVER_PORT},
				{ContainerPort: root.RUSTFS_CONSOLE_PORT},
			},
			VolumeMounts: []core.VolumeMount{
				{
					Name:      dataVol.Name,
					MountPath: mountPath,
				},
				{
					Name:      logsVol.Name,
					MountPath: "/logs",
				},
			},
		}},
		Volumes: []core.Volume{
			dataVol,
			logsVol,
		},
	}
	d := kube.NewDeployment(meta, podSpec)
	// NOTE: For deployments that use PVCs with RWO is better to use the
	// 	recreate strategy in order to avoid "Multi-Attach error for volume" errors.
	d.Spec.Strategy = apps.DeploymentStrategy{Type: apps.RecreateDeploymentStrategyType}
	return d
}
