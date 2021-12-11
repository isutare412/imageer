package config

type Config struct {
	Server ServerConfig `yaml:"server" json:"server"`
	Redis  RedisConfig  `yaml:"redis" json:"redis"`
}

type ServerConfig struct {
	Mode string     `yaml:"mode" json:"mode"`
	Http HttpConfig `yaml:"http" json:"http"`
}

type HttpConfig struct {
	Host string `yaml:"host" json:"host"`
	Port string `yaml:"port" json:"port"`
}

type RedisConfig struct {
	Addrs    []string `yaml:"addrs" json:"addrs"`
	Password string   `yaml:"password" json:"password"`
}
