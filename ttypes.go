package nebula_go_sdk

type HostAddress struct {
	Host string
	Port int
}

type timezoneInfo struct {
	offset int32
	name   []byte
}
