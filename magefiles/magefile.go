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
var Default = Build
var Aliases = map[string]any{
	"b": Build,
}

func init() {
	Env := map[string]string{
		"LINUX_ADMIN":        "arch",
		"CLUSTER_NAME":       root.HYDRA_CLUSTER,
		"CLUSTER_HOST":       root.HYDRA_HOSTNAME,
		"GITEA_HOST":         root.GITEA_HOST,
		"GO_SECRETS":         root.GO_SECRETS_FOLDER,
		"ENC_SECRETS_PATH":   root.GO_ENC_SECRETS_FOLDER,
		"HYDRA_PATH":         root.FLUX_APPS_HYDRA_PATH,
		"SECRETS_FOLDER":     root.SECRETS_FOLDER,
		"HYDRA_SECRETS_PATH": root.FLUX_APPS_SECRETS_HYDRA_PATH,
	}
	r = target.NewRunner(Env, nil)
}

func Build() error {
	t := target.NewA("go", "run", "./cmd/apps/main.go")
	return r.RunV("Build apps", t)
}

func Build_secrets() error {
	t := target.NewA("go", "run", "./cmd/secrets/main.go")
	return r.RunV("Build secrets", t)
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

// Uploads AGE Key to flux as a Kubernetes secret
func Enc_flux() error {
	t := target.New("./scripts/encrypt_flux.sh")
	return r.RunV("Upload key to flux", t)
}

// Encrypt and decrypt Go files that contain secrets.
func Enc_Dec(opt string) error {
	t := target.NewA("./scripts/enc_dec_go.sh", opt)
	return r.RunV("Enc or Dec Go secrets", t)
}
