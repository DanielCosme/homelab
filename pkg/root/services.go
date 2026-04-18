package root

type Service struct {
	Name              string
	Image             string
	Port              int32
	SecurityContextID int64
}

var Linkding = Service{
	Name:              "linkding",
	Image:             "sissbruecker/linkding:1.45.0",
	Port:              9090,
	SecurityContextID: 33, // www-data user, group and FS ID
}
