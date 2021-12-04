package config

type Config struct {
	Server ServerConfig `yaml:"server" json:"server"`
	Redis  RedisConfig  `yaml:"redis" json:"redis"`
}

type ServerConfig struct {
	Mode string `yaml:"mode" json:"mode"`
}

type RedisConfig struct {
	Addrs    []string `yaml:"addrs" json:"addrs"`
	Password string   `yaml:"password" json:"password"`
}
