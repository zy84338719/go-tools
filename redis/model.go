package redis

type Redis struct {
	DB          int
	Addr        string
	ClusterMode bool
	Key         string
	Password    string
	Pattern     string
	BatchSize   int64
	Types       string
}
