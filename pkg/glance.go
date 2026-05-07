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
var config_map core.ConfigMap

func init() {
	meta = kube.NewMetadata(root.Glance.Name, Namespace)
	srv = meta.Service(root.Glance.Port)
	srv.Spec.Type = core.ServiceTypeNodePort
	srv.Spec.Ports[0].NodePort = int32(root.Glance.Public.NodePort)

	config_map = kube.ConfigFromFile("glance.yml", "./config/glance/config/glance.yml", meta)
	config_map.Data["home.yml"] = string(kube.ReadFileBytes("./config/glance/config/home.yml"))
	config_map.Data["user.css"] = string(kube.ReadFileBytes("./config/glance/assets/user.css"))
}

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"namespace": Namespace,
			"service":   srv,
			"configmap": config_map,
			// "deployment": Deployment(),
		},
	)
	return kz.Stack("glance")
}

func Deployment() apps.Deployment {
	// MY_SECRET_TOKEN=123456  ?
	// Config volume - ReadOnly?
	// Assets volume - ReadOnly?
	// Mount /etc/localtime (ReadOnly)
	podSpec := core.PodSpec{
		Containers: []core.Container{{
			Name:  root.Glance.Name,
			Image: root.Glance.Image,
			Env:   kube.NewEnvVarWithSecret(nil, nil, Secret.Name),
		}},
		Volumes: []core.Volume{
			kube.NewVolumeFrom(kube.VolumeSourceConfigMap, "config", config_map.Name),
		},
	}
	return kube.NewDeployment(meta, podSpec)
}
