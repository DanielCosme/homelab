package root

import (
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
	Image: "",
	// PublicURL: "https://photos.danicos.me",
}

var (
	Longhorn           = "longhorn"
	Monitoring         = "monitoring"
	CloudNativePG      = "cnpg"
	PGClusterNamespace = kube.Namespace(CloudNativePG + "-cluster")
)

var (
	TrueNAS_CSI               = "truenas-csi"
	TrueNASURL                = "apex-truenas.orca-uaru.ts.net"
	TrueNASProvisioner        = "cis.truenas.io"
	TrueNASSTorageClassNFS    = "truenas-nfs"
	TrueNASSTorageClass_iSCSI = "truenas-iscsi"
)
