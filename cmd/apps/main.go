package main

import (
	"fmt"
	"os"

	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/flux"
	"danicos.dev/daniel/homelab/pkg/linkding"
	"danicos.dev/daniel/homelab/pkg/monitoring"
	"danicos.dev/daniel/homelab/pkg/root"
	/*
		apps "k8s.io/api/apps/v1"
		core "k8s.io/api/core/v1"
		net "k8s.io/api/networking/v1"
	*/)

func main() {
	flux_stack := flux.Stack()
	err := flux_stack.MarshalYamlFlat(root.FLUX_CLUSTER_HYDRA_PATH)
	assertNoErr(err)

	hydra_monitoring := map[string]stack.Stack{
		"controllers": monitoring.Controllers(),
	}
	for name, s := range hydra_monitoring {
		fmt.Printf("STACK: %s\n", name)
		err = s.MarshalYaml(root.FLUX_INFRA_HYDRA_PATH + "/monitoring")
		assertNoErr(err)
	}

	hydra_apps := map[string]stack.Stack{
		"linkding": linkding.Stack(),
	}
	for name, s := range hydra_apps {
		fmt.Printf("STACK: %s\n", name)
		err = s.MarshalYaml(root.FLUX_APPS_HYDRA_PATH)
		assertNoErr(err)
	}
}

func assertNoErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
