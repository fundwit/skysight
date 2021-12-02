package introspect

type ServiceMeta struct {
	Name string
	ID   string
}

var meta = ServiceMeta{Name: "skysight", ID: "xxx"}

func GetServiceMeta() ServiceMeta {
	return meta
}
