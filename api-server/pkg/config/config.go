package config

type Config struct {
	Server ServerConfig `yaml:"server" json:"server"`
	Redis  RedisConfig  `yaml:"redis" json:"redis"`
	MySQL  MySQLConfig  `yaml:"mysql" json:"mysql"`
	Auth   AuthConfig   `yaml:"auth" json:"auth"`
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

type MySQLConfig struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Address  string `yaml:"address" json:"address"`
	Database string `yaml:"database" json:"database"`
}

type AuthConfig struct {
	ExpireHour int64  `yaml:"expireHour" json:"expireHour"`
	PrivateKey string `yaml:"privateKey" json:"privateKey"`
	PublicKey  string `yaml:"publicKey" json:"publicKey"`
}
