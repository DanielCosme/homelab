package linkding

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

var meta kube.Metadata
var Namespace = kube.Namespace("linkding")
var srv core.Service

func init() {
	meta = kube.NewMetadata("linking", Namespace)
	srv = meta.Service(root.Linkding.Port)
}

func Stack() stack.Stack {
	s := stack.NewStack("linkding", map[string]any{
		"namespace":  Namespace,
		"srv":        srv,
		"deployment": deployment(),
	})
	return s
}

func deployment() apps.Deployment {
	pod_spec := core.PodSpec{
		Containers: []core.Container{
			{
				Name:  root.Linkding.Name,
				Image: root.Linkding.Image,
				Ports: []core.ContainerPort{{
					ContainerPort: root.Linkding.Port,
				}},
			},
		},
	}
	return kube.NewDeployment(meta, pod_spec)
}
