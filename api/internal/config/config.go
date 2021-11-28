package config

type Config struct {
	Server ServerConfig `yaml:"server" json:"server"`
}

type ServerConfig struct {
	Mode string     `yaml:"mode" json:"mode"`
	Http HttpConfig `yaml:"http" json:"http"`
}

type HttpConfig struct {
	Host string `yaml:"host" json:"host"`
	Port string `yaml:"port" json:"port"`
}
