package main

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/homelab/pkg/flux"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
)

func main() {
	meta := kube.NewMetadata("name", kube.Namespace("nae"))
	kuz := kube.NewKustomization(meta, kustomizev1.KustomizationSpec{})
	fluxStack := flux.Stack()
	fluxStack.Add("apps", kuz)
	fluxStack.MarshalYaml("")
}
