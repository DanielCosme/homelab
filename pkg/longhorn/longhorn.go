package longhorn

import (
	"time"

	"danicos.dev/daniel/go-kube/pkg/flux"
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	helm "github.com/fluxcd/helm-controller/api/v2"
	source "github.com/fluxcd/source-controller/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var meta kube.Metadata
var Namespace = kube.Namespace(root.Longhorn + "-system")

func init() {
	meta = kube.NewMetadata(root.Longhorn, Namespace)
}

func Stack() stack.Stack {
	s := stack.NewStack("longhorn", map[string]any{
		"namespace":      Namespace,
		"longhorn-stack": LonghornHelmSource(),
		"release":        LonghornHelmRelease(),
	})
	return s
}

func LonghornHelmSource() source.HelmRepository {
	spec := source.HelmRepositorySpec{
		Interval: durHour(root.FLUX_HELM_MONITORING_INTERVAL),
		URL:      root.HELM_LONGHORN_URL,
	}
	return flux.NewFluxHelmRepositorySource(meta, spec)
}

func LonghornHelmRelease() helm.HelmRelease {
	interval := durHour(12)
	spec := helm.HelmReleaseSpec{
		Interval: durMin(30),
		Chart: &helm.HelmChartTemplate{
			Spec: helm.HelmChartTemplateSpec{
				Chart:    root.HELM_LONGHORN_CHART,
				Version:  root.HELM_LONGHORN_CHART_VERSION,
				Interval: &interval,
				SourceRef: helm.CrossNamespaceObjectReference{
					Kind:      flux.MetaHelmRepository.Kind,
					Name:      meta.Meta().Name,
					Namespace: Namespace.Name,
				},
			},
		},
		Install: &helm.Install{
			CRDs: helm.Create,
		},
		Upgrade: &helm.Upgrade{
			CRDs: helm.CreateReplace,
		},
	}
	return flux.NewFluxHelmRelease(meta, spec)
}

func durHour(d int64) metav1.Duration {
	return metav1.Duration{Duration: (time.Duration(d) * time.Hour)}
}

func durMin(d int64) metav1.Duration {
	return metav1.Duration{Duration: (time.Duration(d) * time.Minute)}
}
