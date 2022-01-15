package config

type Config struct {
	Server ServerConfig `yaml:"server" json:"server"`
	Redis  RedisConfig  `yaml:"redis" json:"redis"`
	MySQL  MySQLConfig  `yaml:"mysql" json:"mysql"`
	S3     S3Config     `yaml:"s3" json:"s3"`
	Auth   AuthConfig   `yaml:"auth" json:"auth"`
}

type ServerConfig struct {
	Mode string     `yaml:"mode" json:"mode"`
	Http HttpConfig `yaml:"http" json:"http"`
	Job  JobConfig  `yaml:"job" json:"job"`
}

type JobConfig struct {
	Queue JobQueueConfig `yaml:"queue" json:"queue"`
	Repo  JobRepoConfig  `yaml:"repo" json:"repo"`
}

type JobQueueConfig struct {
	Request  string `yaml:"request" json:"request"`
	Response string `yaml:"response" json:"response"`
}

type JobRepoConfig struct {
	S3 JobRepoS3Config `yaml:"s3" json:"s3"`
}

type JobRepoS3Config struct {
	Bucket    string `yaml:"bucket" json:"bucket"`
	SourceDir string `yaml:"sourceDir" json:"sourceDir"`
	ResultDir string `yaml:"resultDir" json:"resultDir"`
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

type S3Config struct {
	Address   string `yaml:"address" json:"address"`
	AccessKey string `yaml:"accessKey" json:"accessKey"`
	SecretKey string `yaml:"secretKey" json:"secretKey"`
}

type AuthConfig struct {
	ExpireHour int64  `yaml:"expireHour" json:"expireHour"`
	PrivateKey string `yaml:"privateKey" json:"privateKey"`
	PublicKey  string `yaml:"publicKey" json:"publicKey"`
}
