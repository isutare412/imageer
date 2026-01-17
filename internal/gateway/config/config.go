package config

import (
	"time"

	"github.com/isutare412/imageer/pkg/log"
)

type Config struct {
	Log        LogConfig        `koanf:"log"`
	Web        WebConfig        `koanf:"web"`
	Kubernetes KubernetesConfig `koanf:"kubernetes"`
	Database   DatabaseConfig   `koanf:"database"`
	Valkey     ValkeyConfig     `koanf:"valkey"`
	Auth       AuthConfig       `koanf:"auth"`
	Crypt      CryptConfig      `koanf:"crypt"`
	AWS        AWSConfig        `koanf:"aws"`
	Service    ServiceConfig    `koanf:"service"`
}

type LogConfig struct {
	Format    log.Format `koanf:"format" validate:"validateFn=Validate"`
	Level     log.Level  `koanf:"level" validate:"validateFn=Validate"`
	AddSource bool       `koanf:"add-source"`
	Component string     `koanf:"component" validate:"required"`
}

type WebConfig struct {
	Port              int           `koanf:"port" validate:"required,gt=0,lte=65535"`
	ShowBanner        bool          `koanf:"show-banner"`
	ShowOpenAPIDocs   bool          `koanf:"show-openapi-docs"`
	WriteTimeout      time.Duration `koanf:"write-timeout" validate:"omitempty,gt=0"`
	ReadTimeout       time.Duration `koanf:"read-timeout" validate:"omitempty,gt=0"`
	ReadHeaderTimeout time.Duration `koanf:"read-header-timeout" validate:"omitempty,gt=0"`
	CORS              struct {
		AllowOrigins     string        `koanf:"allow-origins" validate:"required"`
		AllowHeaders     string        `koanf:"allow-headers" validate:"required"`
		AllowMethods     string        `koanf:"allow-methods" validate:"required"`
		AllowCredentials bool          `koanf:"allow-credentials"`
		MaxAge           time.Duration `koanf:"max-age" validate:"required,gt=0"`
	} `koanf:"cors"`
}

type KubernetesConfig struct {
	Enabled bool `koanf:"enabled"`

	Config struct {
		UseInCluster bool `koanf:"use-in-cluster"`
		Kubeconfig   struct {
			Path    string `koanf:"path"`
			Context string `koanf:"context"`
		} `koanf:"kubeconfig"`
	} `koanf:"config"`

	LeaderElection struct {
		LeaseName      string        `koanf:"lease-name" validate:"required"`
		LeaseNamespace string        `koanf:"lease-namespace" validate:"required"`
		LeaseDuration  time.Duration `koanf:"lease-duration" validate:"required,gt=0"`
		RenewDeadline  time.Duration `koanf:"renew-deadline" validate:"required,gt=0"`
		RetryPeriod    time.Duration `koanf:"retry-period" validate:"required,gt=0"`
	} `koanf:"leader-election"`
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

type ValkeyConfig struct {
	Addresses string `koanf:"addresses" validate:"required"`
	Username  string `koanf:"username" validate:"required"`
	Password  string `koanf:"password" validate:"required"`

	Streams struct {
		ImageProcessRequest struct {
			StreamKey  string `koanf:"stream-key" validate:"required"`
			StreamSize int    `koanf:"stream-size" validate:"required,gt=0"`
		} `koanf:"image-process-request"`

		ImageProcessResult struct {
			StreamKey string `koanf:"stream-key" validate:"required"`
			GroupName string `koanf:"group-name" validate:"required"`
			Handler   struct {
				Concurrency int           `koanf:"concurrency" validate:"required,gt=0"`
				Timeout     time.Duration `koanf:"timeout" validate:"required,gt=0"`
			} `koanf:"handler"`
			Reader struct {
				BatchSize    int64         `koanf:"batch-size" validate:"required,gt=0"`
				BlockTimeout time.Duration `koanf:"block-timeout" validate:"required,gt=0"`
			} `koanf:"reader"`
			Stealer struct {
				Interval           time.Duration `koanf:"interval" validate:"required,gt=0"`
				MinIdleTime        time.Duration `koanf:"min-idle-time" validate:"required,gt=0"`
				MaxDeliveryAttempt int64         `koanf:"max-delivery-attempt" validate:"required,gt=0"`
			} `koanf:"stealer"`
			Reaper struct {
				MinIdleTime time.Duration `koanf:"min-idle-time" validate:"required,gt=0"`
			} `koanf:"reaper"`
		} `koanf:"image-process-result"`
	} `koanf:"streams"`

	PubSub struct {
		ImageProcessDone struct {
			ChannelPrefix string `koanf:"channel-prefix" validate:"required"`
			MaxRetries    int    `koanf:"max-retries" validate:"gte=0"`
		} `koanf:"image-process-done"`
	} `koanf:"pubsub"`
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

type AWSConfig struct {
	S3         S3Config         `koanf:"s3"`
	CloudFront CloudFrontConfig `koanf:"cloudfront"`
	SQS        SQSConfig        `koanf:"sqs"`
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

type CloudFrontConfig struct {
	Images struct {
		DistributionDomain string `koanf:"distribution-domain" validate:"required"`
	} `koanf:"images"`
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

type ServiceConfig struct {
	Image struct {
		ProcessDoneWaitTimeout time.Duration `koanf:"process-done-wait-timeout" validate:"required,gt=0"`
		ExpireCheckInterval    time.Duration `koanf:"expire-check-interval" validate:"required,gt=0"`
	} `koanf:"image"`
}
