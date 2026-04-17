package main

import (
	"danicos.dev/daniel/homelab/pkg/root"
	"danicos.dev/daniel/homelab/pkg/target"
	// "github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	/*
		NOTE: Mage https://github.com/magefile/mage
		Mage other packages
		"github.com/magefile/mage/mage"
		"github.com/magefile/mage/parse"
		"github.com/magefile/mage/target"
	*/)

var r target.Runner
var Default = Run
var Aliases = map[string]any{
	"r": Run,
}

func init() {
	Env := map[string]string{
		"LINUX_ADMIN":  "arch",
		"CLUSTER_NAME": root.HYDRA_CLUSTER,
		"CLUSTER_HOST": root.HYDRA_HOSTNAME,
		"GITEA_HOST":   root.GITEA_HOST,
	}
	r = target.NewRunner(Env, nil)
}

func Run() error {
	t := target.NewA("go", "run", ".")
	return r.RunV("run", t)
}

func InstallK3S() error {
	t := target.New("./scripts/install_k3s.sh")
	return r.RunV("Install k3s", t)
}

func GetKubeconfig() error {
	t := target.New("./scripts/get_kubeconfig.sh")
	return r.RunV("Get Kubeconfig", t)
}

func InstallFlux() error {
	t := target.New("./scripts/install_flux.sh")
	return r.RunV("Install Flux", t)
}
