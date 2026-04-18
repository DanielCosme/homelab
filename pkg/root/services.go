package root

type Service struct {
	Name  string
	Image string
	Port  int32
}

var Linkding = Service{
	Name:  "linkding",
	Image: "sissbruecket/linkding:1.31.0",
	Port:  9090,
}
