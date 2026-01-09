package config

import (
	"time"

	"github.com/isutare412/imageer/pkg/log"
)

type Config struct {
	Log      LogConfig      `koanf:"log"`
	Web      WebConfig      `koanf:"web"`
	Database DatabaseConfig `koanf:"database"`
	Auth     AuthConfig     `koanf:"auth"`
	Crypt    CryptConfig    `koanf:"crypt"`
	S3       S3Config       `koanf:"s3"`
	SQS      SQSConfig      `koanf:"sqs"`
}

type LogConfig struct {
	Format    log.Format `koanf:"format" validate:"validateFn=Validate"`
	Level     log.Level  `koanf:"level" validate:"validateFn=Validate"`
	AddSource bool       `koanf:"add-source"`
	Component string     `koanf:"component" validate:"required"`
}

type WebConfig struct {
	Port            int  `koanf:"port" validate:"required,gt=0,lte=65535"`
	ShowBanner      bool `koanf:"show-banner"`
	ShowOpenAPIDocs bool `koanf:"show-openapi-docs"`
}

type DatabaseConfig struct {
	UseInMemory bool           `koanf:"use-in-memory"`
	TraceLog    bool           `koanf:"trace-log"`
	Postgres    PostgresConfig `koanf:"postgres"`
}

type PostgresConfig struct {
	Host     string `koanf:"host" validate:"required"`
	Port     int    `koanf:"port" validate:"required,gt=0,lte=65535"`
	User     string `koanf:"user" validate:"required"`
	Password string `koanf:"password" validate:"required"`
	Database string `koanf:"database" validate:"required"`
}

type AuthConfig struct {
	Cookies struct {
		OIDCState struct {
			Name string        `koanf:"name" validate:"required"`
			TTL  time.Duration `koanf:"ttl" validate:"required,gte=1m"`
		} `koanf:"oidc-state"`

		User struct {
			Name string        `koanf:"name" validate:"required"`
			TTL  time.Duration `koanf:"ttl" validate:"required,gte=1m"`
		} `koanf:"user"`
	} `koanf:"cookies"`

	ServiceAccount struct {
		APIKeyHeader string `koanf:"api-key-header" validate:"required"`
	} `koanf:"service-account"`

	JWT struct {
		ActiveKeyPairName string                       `koanf:"active-key-pair-name" validate:"required"`
		KeyPairs          map[string]AuthKeyPairConfig `koanf:"key-pairs" validate:"required,dive,keys,required,endkeys,required"`
	} `koanf:"jwt"`

	Google struct {
		ClientID     string `koanf:"client-id" validate:"required"`
		ClientSecret string `koanf:"client-secret" validate:"required"`
		CallbackPath string `koanf:"callback-path" validate:"required"`
	} `koanf:"google"`
}

type AuthKeyPairConfig struct {
	Private string `koanf:"private" validate:"required"`
	Public  string `koanf:"public" validate:"required"`
}

type CryptConfig struct {
	AES struct {
		Key string `koanf:"key" validate:"required"`
	} `koanf:"aes"`
}

type S3Config struct {
	Bucket string `koanf:"bucket" validate:"required"`
	Prefix struct {
		Image string `koanf:"image" validate:"required"`
	} `koanf:"prefix"`
	Presign struct {
		Expiry time.Duration `koanf:"expiry" validate:"required,gt=0"`
	} `koanf:"presign"`
}

type SQSConfig struct {
	ImageUploadEventQueue struct {
		QueueURL           string        `koanf:"queue-url" validate:"required"`
		BatchCount         int32         `koanf:"batch-count" validate:"required,gt=0,lte=10"`
		PollingWaitTimeout time.Duration `koanf:"polling-wait-timeout" validate:"required,gt=0,lte=20s"`
		VisibilityTimeout  time.Duration `koanf:"visibility-timeout" validate:"required,gt=0,lte=12h"`
		BatchHandleTimeout time.Duration `koanf:"batch-handle-timeout" validate:"required,gt=0,lte=10m"`
		HandleTimeout      time.Duration `koanf:"handle-timeout" validate:"required,gt=0,lte=10m"`
	} `koanf:"image-upload-event-queue"`
}
