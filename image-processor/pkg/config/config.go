package config

type Config struct {
	Server    ServerConfig    `yaml:"server" json:"server"`
	Redis     RedisConfig     `yaml:"redis" json:"redis"`
	Processor ProcessorConfig `yaml:"processor" json:"processor"`
}

type ServerConfig struct {
	Mode string `yaml:"mode" json:"mode"`
}

type RedisConfig struct {
	Addrs    []string          `yaml:"addrs" json:"addrs"`
	Password string            `yaml:"password" json:"password"`
	Stream   RedisStreamConfig `yaml:"stream" json:"stream"`
}

type RedisStreamConfig struct {
	GroupName string `yaml:"groupName" json:"groupName"`
}

type ProcessorConfig struct {
	RetryDelay int64 `yaml:"retryDelay" json:"retryDelay"`
}
