package flux

import (
	"time"

	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	kz "github.com/fluxcd/kustomize-controller/api/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var flux_apps_meta kube.Metadata
var Flux_namespace = kube.Namespace(root.FLUX_NAMESPACE)

func init() {
	flux_apps_meta = kube.NewMetadata("apps", Flux_namespace)
}

func Stack() stack.Stack {
	s := stack.NewStack("flux", map[string]any{
		"apps": Apps(),
	})
	return s
}

func Apps() kz.Kustomization {
	retryInteval := durMin(1)
	timeout := durMin(5)
	spec := kz.KustomizationSpec{
		Interval:      durMin(10),
		RetryInterval: &retryInteval,
		Timeout:       &timeout,
		SourceRef: kz.CrossNamespaceSourceReference{
			Kind: "GitRepository",
			Name: Flux_namespace.Name,
		},
		Path:  root.FLUX_APPS_HYDRA_PATH,
		Prune: true,
	}
	return kube.NewKustomization(flux_apps_meta, spec)
}

func durMin(d int64) meta.Duration {
	return meta.Duration{Duration: (time.Duration(d) * time.Minute)}
}
