package vaultwarden

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
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
			"namespace": Namespace,
			"service":   srv,
			"pvc":       pvc,
			// "deployment":       Deployment(),
		},
	)
	return kz.Stack("vaultwarden")
}

/*
services:
  vaultwarden:
    image: vaultwarden/server:latest
    container_name: vaultwarden
    restart: always
    environment:
      # DOMAIN: "https://vaultwarden.example.com"  # required when using a reverse proxy; your domain; vaultwarden needs to know it's https to work properly with attachments
      SIGNUPS_ALLOWED: "true" # Deactivate this with "false" after you have created your account so that no strangers can register
    volumes:
      - ./vw-data:/data # the path before the : can be changed
    ports:
      - 11001:80 # you can replace the 11001 with your preferred port

*/
