package root

type Service struct {
	Name              string
	Image             string
	PublicURL         string
	Port              int32
	SecurityContextID int64
}

var Linkding = Service{
	Name:              "linkding",
	Image:             "sissbruecker/linkding:1.45.0",
	PublicURL:         "https://link.danicos.me",
	Port:              9090,
	SecurityContextID: 33, // www-data user, group and FS ID
}

var (
	Longhorn    = "longhorn"
	Monitoring  = "monitoring"
	TrueNAS_CSI = "truenas-csi"
	TrueNASURL  = "apex-truenas.orca-uaru.ts.net"
)
