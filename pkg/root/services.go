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
	Image: "ghcr.io/immich-app/immich-server:v2.6.0:-release}",
	Port:  2283,
	// PublicURL: "https://photos.danicos.me",
}

var (
	Longhorn            = "longhorn"
	Monitoring          = "monitoring"
	CloudNativePG       = "cnpg"
	PGClusterNamespace  = kube.Namespace(CloudNativePG + "-cluster")
	PG_CLUSTER_HOSTNAME = fmt.Sprintf("%s-rw.%s.svc.cluster.local", PGClusterNamespace.Name, PGClusterNamespace.Name)
	PG_CLUSTER_PORT     = 5432
)

var (
	TrueNAS_CSI               = "truenas-csi"
	TrueNASURL                = "apex-truenas.orca-uaru.ts.net"
	TrueNASProvisioner        = "cis.truenas.io"
	TrueNASSTorageClassNFS    = "truenas-nfs"
	TrueNASSTorageClass_iSCSI = "truenas-iscsi"
)
