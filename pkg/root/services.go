package root

import (
	"fmt"

	"danicos.dev/daniel/go-kube/pkg/kube"
)

type Service struct {
	Name              string
	Image             string
	Port              int32
	SecurityContextID int64
	Public            *Public
	Postgres          *Postgres
}

type Public struct {
	URL      string
	NodePort int32
}

type Postgres struct{}

var Linkding = Service{
	Name:              "linking",
	Image:             "sissbruecker/linkding:1.45.0",
	Port:              9090,
	SecurityContextID: 33, // www-data user, group and FS ID
	Public: &Public{
		URL:      "https://link.danicos.me",
		NodePort: 30010,
	},
}

var Immich = Service{
	Name:  "immich",
	Image: "ghcr.io/immich-app/immich-server:v2.7.5",
	Port:  2283, // Server Port
	Public: &Public{
		URL:      "https://photos.danicos.me",
		NodePort: 30011,
	},
}

var Glance = Service{
	Name:  "glance",
	Image: "glanceapp/glance:v0.8.4",
	Port:  8080, // Server Port
	Public: &Public{
		URL:      "https://home.danicos.me",
		NodePort: 30009,
	},
}

var Vaultwarden = Service{
	Name:  "vaultwarden",
	Image: "quay.io/vaultwarden/server:1.37.4",
	Port:  80, // Server Port
	Public: &Public{
		URL:      "https://vault.danicos.me",
		NodePort: 30012,
	},
}

var (
	Longhorn            = "longhorn"
	Monitoring          = "monitoring"
	CloudNativePG       = "cnpg"
	PGClusterNamespace  = kube.Namespace(CloudNativePG + "-cluster")
	PG_CLUSTER_HOSTNAME = fmt.Sprintf("%s-rw.%s.svc.cluster.local", PGClusterNamespace.Name, PGClusterNamespace.Name)
	PG_CLUSTER_PORT     = 5432
)
