package config

type Config struct {
	Mode  string      `yaml:"mode" json:"mode"`
	Job   JobConfig   `yaml:"job" json:"job"`
	Redis RedisConfig `yaml:"redis" json:"redis"`
}

type JobConfig struct {
	RetryDelay int64          `yaml:"retryDelay" json:"retryDelay"`
	Queue      JobQueueConfig `yaml:"queue" json:"queue"`
}

type JobQueueConfig struct {
	Request  string `yaml:"request" json:"request"`
	Response string `yaml:"response" json:"response"`
}

type RedisConfig struct {
	Addrs    []string          `yaml:"addrs" json:"addrs"`
	Password string            `yaml:"password" json:"password"`
	Stream   RedisStreamConfig `yaml:"stream" json:"stream"`
}

type RedisStreamConfig struct {
	GroupName string `yaml:"groupName" json:"groupName"`
}
