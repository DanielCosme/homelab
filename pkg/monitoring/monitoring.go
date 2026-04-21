package monitoring

import (
	"encoding/json"
	"time"

	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	helm "github.com/fluxcd/helm-controller/api/v2"
	"github.com/fluxcd/pkg/apis/kustomize"
	source "github.com/fluxcd/source-controller/api/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var meta kube.Metadata
var Namespace = kube.Namespace(root.Monitoring)

func init() {
	meta = kube.NewMetadata(root.Monitoring, Namespace)
}

func Controllers() stack.Stack {
	s := stack.NewStack("controllers", map[string]any{
		"namespace":             Namespace,
		"kube-prometheus-stack": PrometheusHelmSource(),
		"release":               PrometheusRelease(),
	})

	return s
}

func PrometheusHelmSource() source.HelmRepository {
	spec := source.HelmRepositorySpec{
		Interval: durHour(root.FLUX_HELM_MONITORING_INTERVAL),
		URL:      root.HELM_PROMETHEUS_URL,
	}
	return kube.NewFluxHelmRepositorySource(meta, spec)
}

func PrometheusRelease() helm.HelmRelease {
	type Grafana struct {
		AdminPassword string `json:"adminPassword"`
	}
	type Values struct {
		Grafana Grafana `json:"grafana"`
	}
	values := Values{
		Grafana: Grafana{
			// Note this password is safe because Grafana is only exposed via VPN
			AdminPassword: "grafana-admin",
		},
	}
	raw, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}
	interval := durHour(12)
	spec := helm.HelmReleaseSpec{
		Interval: durMin(30),
		Chart: &helm.HelmChartTemplate{
			Spec: helm.HelmChartTemplateSpec{
				Chart:    root.HELM_PROMETHEUS_CHART,
				Version:  root.HELM_PROMETHEUS_CHART_VERSION,
				Interval: &interval,
				SourceRef: helm.CrossNamespaceObjectReference{
					Kind:      kube.FluxHelmReleaseMeta.Kind,
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
		DriftDetection: &helm.DriftDetection{
			Mode: helm.DriftDetectionEnabled,
			Ignore: []helm.IgnoreRule{{
				Paths: []string{"/metadata/annotations/prometheus-operator-validated"},
				Target: &kustomize.Selector{
					Kind: "PrometheusRule",
				},
			}},
		},
		Values: &apiextensionsv1.JSON{Raw: raw},
	}
	return kube.NewFluxHelmRelease(meta, spec)
}

func durHour(d int64) metav1.Duration {
	return metav1.Duration{Duration: (time.Duration(d) * time.Hour)}
}

func durMin(d int64) metav1.Duration {
	return metav1.Duration{Duration: (time.Duration(d) * time.Minute)}
}
