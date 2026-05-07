package glance

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

var Secret = struct {
	Name     string
	TokenKey string
}{
	Name:     root.Glance.Name + "-secret",
	TokenKey: "token",
}

var Namespace = kube.Namespace(root.Glance.Name)
var meta kube.Metadata
var srv core.Service
var config_glance core.ConfigMap
var config_assets core.ConfigMap

func init() {
	meta = kube.NewMetadata(root.Glance.Name, Namespace)
	srv = meta.Service(root.Glance.Port)
	srv.Spec.Type = core.ServiceTypeNodePort
	srv.Spec.Ports[0].NodePort = int32(root.Glance.Public.NodePort)

	config_glance = kube.ConfigFromFile("glance.yml", "./config/glance/config/glance.yml", meta)
	config_glance.Data["home.yml"] = string(kube.ReadFileBytes("./config/glance/config/home.yml"))
	meta_assets := kube.NewMetadata(root.Glance.Name+"-assets", Namespace)
	config_assets = kube.ConfigFromFile("user.css", "./config/glance/assets/user.css", meta_assets)
}

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"namespace":        Namespace,
			"service":          srv,
			"configmap-glance": config_glance,
			"configmap-assets": config_assets,
			"deployment":       Deployment(),
		},
	)
	return kz.Stack("glance")
}

func Deployment() apps.Deployment {
	configVol := kube.NewVolumeFrom(kube.VolumeSourceConfigMap, "config", config_glance.Name)
	assetsVol := kube.NewVolumeFrom(kube.VolumeSourceConfigMap, "assets", config_assets.Name)
	localtimeVol := core.Volume{
		Name: "localtime",
		VolumeSource: core.VolumeSource{
			HostPath: &core.HostPathVolumeSource{
				Path: "/etc/localtime",
				Type: new(core.HostPathFile),
			},
		},
	}
	podSpec := core.PodSpec{
		Containers: []core.Container{{
			Name:  root.Glance.Name,
			Image: root.Glance.Image,
			// Env:   kube.NewEnvVarWithSecret(nil, nil, Secret.Name),
			VolumeMounts: []core.VolumeMount{
				{
					Name:      configVol.Name,
					MountPath: "/app/config",
				},
				{
					Name:      assetsVol.Name,
					MountPath: "/app/assets",
				},
				{
					Name:      localtimeVol.Name,
					MountPath: "/etc/localtime",
					ReadOnly:  true,
				},
			},
		}},
		Volumes: []core.Volume{
			configVol,
			assetsVol,
			localtimeVol,
		},
	}
	return kube.NewDeployment(meta, podSpec)
}
