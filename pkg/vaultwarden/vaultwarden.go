package vaultwarden

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

var Namespace = kube.Namespace(root.Vaultwarden.Name)
var meta kube.Metadata
var srv core.Service
var pvc core.PersistentVolumeClaim

func init() {
	meta = kube.NewMetadata(root.Vaultwarden.Name, Namespace)
	srv = meta.Service(root.Vaultwarden.Port)
	srv.Spec.Type = core.ServiceTypeNodePort
	srv.Spec.Ports[0].NodePort = int32(root.Vaultwarden.Public.NodePort)

	pvc = meta.PVC()
}

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"namespace":  Namespace,
			"service":    srv,
			"pvc":        pvc,
			"deployment": Deployment(),
		},
	)
	return kz.Stack("vaultwarden")
}

func Deployment() apps.Deployment {
	envMap := map[string]string{
		"DOMAIN": root.Vaultwarden.Public.URL,
		// "SIGNUPS_ALLOWED": "true",
	}
	dataVol := kube.NewVolumeFrom(kube.VolumeSourcePVC, "data", pvc.Name)
	podSpec := core.PodSpec{
		Containers: []core.Container{{
			Name:  root.Vaultwarden.Name,
			Image: root.Vaultwarden.Image,
			Ports: []core.ContainerPort{{ContainerPort: root.Vaultwarden.Port}},
			Env:   kube.NewEnvVar(envMap),
			VolumeMounts: []core.VolumeMount{{
				Name:      dataVol.Name,
				MountPath: "/data",
			}},
		}},
		Volumes: []core.Volume{
			dataVol,
		},
	}
	d := kube.NewDeployment(meta, podSpec)
	// NOTE: For deployments that use PVCs with RWO (from longhorn) is better to use the
	// 	recreate strategy in order to avoid "Multi-Attach error for volume" errors.
	d.Spec.Strategy = apps.DeploymentStrategy{Type: apps.RecreateDeploymentStrategyType}
	return d
}
