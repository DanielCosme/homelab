package linkding

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

var Secret = struct {
	Name                 string
	SuperUserKey         string
	SuperUserPasswordKey string
}{
	Name:                 root.Linkding.Name,
	SuperUserKey:         "supe_user_name",
	SuperUserPasswordKey: "supe_user_password",
}

var meta kube.Metadata
var Namespace = kube.Namespace("linkding")
var srv core.Service
var pvc core.PersistentVolumeClaim

func init() {
	meta = kube.NewMetadata("linking", Namespace)
	srv = meta.Service(root.Linkding.Port)
	srv.Spec.Type = core.ServiceTypeNodePort
	srv.Spec.Ports[0].NodePort = 30010
	pvc = meta.PVC()
}

func Stack() stack.Stack {
	s := stack.NewStack("linkding", map[string]any{
		"namespace":  Namespace,
		"srv":        srv,
		"pvc":        pvc,
		"deployment": deployment(),
	})
	return s
}

func deployment() apps.Deployment {
	storage := kube.NewVolumeFrom(kube.VolumeSourcePVC, "data", pvc.Name)
	envMapping := map[string]string{
		"LD_CSRF_TRUSTED_ORIGINS": "https://link.danicos.me",
	}
	secretMapping := map[string]string{
		"LD_SUPERUSER_NAME":     Secret.SuperUserKey,
		"LD_SUPERUSER_PASSWORD": Secret.SuperUserPasswordKey,
	}
	pod_spec := core.PodSpec{
		SecurityContext: &core.PodSecurityContext{
			RunAsUser:  &root.Linkding.SecurityContextID,
			RunAsGroup: &root.Linkding.SecurityContextID,
			FSGroup:    &root.Linkding.SecurityContextID,
		},
		Containers: []core.Container{
			{
				Name:            root.Linkding.Name,
				Image:           root.Linkding.Image,
				SecurityContext: root.ContainerSecurityContext,
				Env:             kube.NewEnvVarWithSecret(envMapping, secretMapping, Secret.Name),
				Ports: []core.ContainerPort{{
					ContainerPort: root.Linkding.Port,
				}},
				VolumeMounts: []core.VolumeMount{{
					Name:      storage.Name,
					MountPath: "/etc/linkding/data",
				}},
			},
		},
		Volumes: []core.Volume{
			storage,
		},
	}
	return kube.NewDeployment(meta, pod_spec)
}

// kubectl -n linkding exec -it linking-67f686679d-2tfrk -- python manage.py createsuperuser --username=daniel --email=danicosme@pm.me
