package config

import (
	"strings"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/crypt"
	"github.com/isutare412/imageer/internal/gateway/jwt"
	"github.com/isutare412/imageer/internal/gateway/kafka"
	"github.com/isutare412/imageer/internal/gateway/kubernetes"
	"github.com/isutare412/imageer/internal/gateway/oidc"
	"github.com/isutare412/imageer/internal/gateway/postgres"
	"github.com/isutare412/imageer/internal/gateway/s3"
	"github.com/isutare412/imageer/internal/gateway/service/auth"
	"github.com/isutare412/imageer/internal/gateway/service/image"
	"github.com/isutare412/imageer/internal/gateway/sqs"
	"github.com/isutare412/imageer/internal/gateway/valkey"
	"github.com/isutare412/imageer/internal/gateway/web"
	"github.com/isutare412/imageer/internal/gateway/webv2"
	"github.com/isutare412/imageer/pkg/log"
	"github.com/isutare412/imageer/pkg/tracing"
)

func (c *Config) ToLogConfig() log.Config {
	return log.Config{
		Format:    c.Log.Format,
		Level:     c.Log.Level,
		AddSource: c.Log.AddSource,
		Component: c.Log.Component,
	}
}

func (c *Config) ToTracingConfig() tracing.Config {
	return tracing.Config(c.Trace)
}

func (c *Config) ToWebConfig() web.Config {
	return web.Config{
		Port:              c.Web.Port,
		ShowBanner:        c.Web.ShowBanner,
		ShowOpenAPIDocs:   c.Web.ShowOpenAPIDocs,
		APIKeyHeader:      c.Auth.ServiceAccount.APIKeyHeader,
		UserCookieName:    c.Auth.Cookies.User.Name,
		WriteTimeout:      c.Web.WriteTimeout,
		ReadTimeout:       c.Web.ReadTimeout,
		ReadHeaderTimeout: c.Web.ReadHeaderTimeout,
		CORS: web.CORSConfig{
			AllowOrigins:     parseCSV(c.Web.CORS.AllowOrigins, ","),
			AllowHeaders:     parseCSV(c.Web.CORS.AllowHeaders, ","),
			AllowMethods:     parseCSV(c.Web.CORS.AllowMethods, ","),
			AllowCredentials: c.Web.CORS.AllowCredentials,
			MaxAge:           c.Web.CORS.MaxAge,
		},
	}
}

func (c *Config) ToWebV2Config() webv2.Config {
	return webv2.Config{
		Port:                  c.Web.Port,
		ShowOpenAPIDocs:       c.Web.ShowOpenAPIDocs,
		APIKeyHeader:          c.Auth.ServiceAccount.APIKeyHeader,
		UserCookieName:        c.Auth.Cookies.User.Name,
		TokenRefreshThreshold: c.Auth.Cookies.User.RefreshThreshold,
		WriteTimeout:          c.Web.WriteTimeout,
		ReadTimeout:           c.Web.ReadTimeout,
		ReadHeaderTimeout:     c.Web.ReadHeaderTimeout,
		CORS: webv2.CORSConfig{
			AllowOrigins:     parseCSV(c.Web.CORS.AllowOrigins, ","),
			AllowHeaders:     parseCSV(c.Web.CORS.AllowHeaders, ","),
			AllowMethods:     parseCSV(c.Web.CORS.AllowMethods, ","),
			AllowCredentials: c.Web.CORS.AllowCredentials,
			MaxAge:           c.Web.CORS.MaxAge,
		},
	}
}

func (c *Config) ToKubernetesClientConfig() kubernetes.ClientConfig {
	return kubernetes.ClientConfig{
		UseInClusterConfig: c.Kubernetes.Config.UseInCluster,
		KubeConfigPath:     c.Kubernetes.Config.Kubeconfig.Path,
		KubeConfigContext:  c.Kubernetes.Config.Kubeconfig.Context,
	}
}

func (c *Config) ToKubernetesLeaderElectorConfig() kubernetes.LeaderElectorConfig {
	return kubernetes.LeaderElectorConfig{
		LeaseName:      c.Kubernetes.LeaderElection.LeaseName,
		LeaseNamespace: c.Kubernetes.LeaderElection.LeaseNamespace,
		LeaseDuration:  c.Kubernetes.LeaderElection.LeaseDuration,
		RenewDeadline:  c.Kubernetes.LeaderElection.RenewDeadline,
		RetryPeriod:    c.Kubernetes.LeaderElection.RetryPeriod,
	}
}

func (c *Config) ToRepositoryClientConfig() postgres.ClientConfig {
	return postgres.ClientConfig{
		Host:        c.Database.Postgres.Host,
		Port:        c.Database.Postgres.Port,
		User:        c.Database.Postgres.User,
		Password:    c.Database.Postgres.Password,
		Database:    c.Database.Postgres.Database,
		TraceLog:    c.Database.TraceLog,
		UseInMemory: c.Database.UseInMemory,
	}
}

func (c *Config) ToValkeyClientConfig() valkey.ClientConfig {
	return valkey.ClientConfig{
		Addresses: parseCSV(c.Valkey.Addresses, ","),
		Username:  c.Valkey.Username,
		Password:  c.Valkey.Password,
	}
}

func (c *Config) ToValkeyImageNotificationPublisherConfig() valkey.ImageNotificationPublisherConfig {
	return valkey.ImageNotificationPublisherConfig{
		UploadDoneChannelPrefix:  c.Valkey.PubSub.ImageUploadDone.ChannelPrefix,
		ProcessDoneChannelPrefix: c.Valkey.PubSub.ImageProcessDone.ChannelPrefix,
	}
}

func (c *Config) ToValkeyImageUploadDoneSubscriberConfig() valkey.ImageUploadDoneSubscriberConfig {
	return valkey.ImageUploadDoneSubscriberConfig{
		ChannelPrefix: c.Valkey.PubSub.ImageUploadDone.ChannelPrefix,
		MaxRetries:    c.Valkey.PubSub.ImageUploadDone.MaxRetries,
	}
}

func (c *Config) ToValkeyImageProcessDoneSubscriberConfig() valkey.ImageProcessDoneSubscriberConfig {
	return valkey.ImageProcessDoneSubscriberConfig{
		ChannelPrefix: c.Valkey.PubSub.ImageProcessDone.ChannelPrefix,
		MaxRetries:    c.Valkey.PubSub.ImageProcessDone.MaxRetries,
	}
}

func (c *Config) ToKafkaClientConfig() kafka.ClientConfig {
	return kafka.ClientConfig{
		Addrs:         parseCSV(c.Kafka.Addresses, ","),
		User:          c.Kafka.Username,
		Password:      c.Kafka.Password,
		ConsumerGroup: c.Kafka.ConsumerGroup,
		ConsumeTopics: []string{
			c.Kafka.Topics.ImageProcessResult.Topic,
			c.Kafka.Topics.ImageProcessResult.RetryTopic,
			c.Kafka.Topics.ImageS3DeleteRequest.Topic,
			c.Kafka.Topics.ImageS3DeleteRequest.RetryTopic,
		},
	}
}

func (c *Config) ToKafkaImageProcessRequestQueueConfig() kafka.ImageProcessRequestQueueConfig {
	return kafka.ImageProcessRequestQueueConfig{
		Topic: c.Kafka.Topics.ImageProcessRequest.Topic,
	}
}

func (c *Config) ToKafkaImageProcessResultHandlerConfig() kafka.ImageProcessResultHandlerConfig {
	return kafka.ImageProcessResultHandlerConfig{
		RetryTopic:      c.Kafka.Topics.ImageProcessResult.RetryTopic,
		HandleTimeout:   c.Kafka.Topics.ImageProcessResult.Handler.Timeout,
		MaxRetryAttempt: c.Kafka.Topics.ImageProcessResult.Handler.MaxRetryAttempt,
		RetryBaseDelay:  c.Kafka.Topics.ImageProcessResult.Handler.RetryBaseDelay,
	}
}

func (c *Config) ToKafkaImageS3DeleteRequestQueueConfig() kafka.ImageS3DeleteRequestQueueConfig {
	return kafka.ImageS3DeleteRequestQueueConfig{
		Topic: c.Kafka.Topics.ImageS3DeleteRequest.Topic,
	}
}

func (c *Config) ToKafkaImageS3DeleteRequestHandlerConfig() kafka.ImageS3DeleteRequestHandlerConfig {
	return kafka.ImageS3DeleteRequestHandlerConfig{
		RetryTopic:      c.Kafka.Topics.ImageS3DeleteRequest.RetryTopic,
		HandleTimeout:   c.Kafka.Topics.ImageS3DeleteRequest.Handler.Timeout,
		MaxRetryAttempt: c.Kafka.Topics.ImageS3DeleteRequest.Handler.MaxRetryAttempt,
		RetryBaseDelay:  c.Kafka.Topics.ImageS3DeleteRequest.Handler.RetryBaseDelay,
	}
}

func (c *Config) ToAuthServiceConfig() auth.ServiceConfig {
	return auth.ServiceConfig{
		StateCookieName: c.Auth.Cookies.OIDCState.Name,
		StateCookieTTL:  c.Auth.Cookies.OIDCState.TTL,
		UserCookieName:  c.Auth.Cookies.User.Name,
		UserCookieTTL:   c.Auth.Cookies.User.TTL,
	}
}

func (c *Config) toRSAKeyPairs() []jwt.RSAKeyBytesPair {
	keyPairs := make([]jwt.RSAKeyBytesPair, 0, len(c.Auth.JWT.KeyPairs))
	for name, pair := range c.Auth.JWT.KeyPairs {
		keyPairs = append(keyPairs, jwt.RSAKeyBytesPair{
			Name:    name,
			Private: []byte(pair.Private),
			Public:  []byte(pair.Public),
		})
	}
	return keyPairs
}

func (c *Config) ToJWTSignerConfig() jwt.SignerConfig {
	return jwt.SignerConfig{
		ActiveKeyPairName: c.Auth.JWT.ActiveKeyPairName,
		KeyPairs:          c.toRSAKeyPairs(),
	}
}

func (c *Config) ToJWTVerifierConfig() jwt.VerifierConfig {
	return jwt.VerifierConfig{
		KeyPairs: c.toRSAKeyPairs(),
	}
}

func (c *Config) ToOIDCGoogleClientConfig() oidc.GoogleClientConfig {
	return oidc.GoogleClientConfig(c.Auth.Google)
}

func (c *Config) ToAESCrypterConfig() crypt.AESCrypterConfig {
	return crypt.AESCrypterConfig(c.Crypt.AES)
}

func (c *Config) ToS3PresignerConfig() s3.PresignerConfig {
	return s3.PresignerConfig{
		Bucket: c.AWS.S3.Bucket,
		Expiry: c.AWS.S3.Presign.Expiry,
	}
}

func (c *Config) ToS3ObjectStorageConfig() s3.ObjectStorageConfig {
	return s3.ObjectStorageConfig{
		Bucket: c.AWS.S3.Bucket,
	}
}

func (c *Config) ToSQSImageUploadListenerConfig() sqs.ImageUploadListenerConfig {
	return sqs.ImageUploadListenerConfig(c.AWS.SQS.ImageUploadEventQueue)
}

func (c *Config) ToImageServiceConfig() image.Config {
	return image.Config{
		CDNDomain:              c.AWS.CloudFront.Images.DistributionDomain,
		S3KeyPrefix:            c.AWS.S3.Prefix.Image,
		ProcessDoneWaitTimeout: c.Service.Image.ProcessDoneWaitTimeout,
	}
}

func (c *Config) ToImageCloserConfig() image.CloserConfig {
	return image.CloserConfig{
		CheckInterval:  c.Service.Image.ExpireCheckInterval,
		CloseThreshold: c.AWS.S3.Presign.Expiry,
	}
}

func parseCSV(s string, delim string) []string {
	parts := strings.Split(s, delim)
	parts = lo.Map(parts, func(item string, _ int) string { return strings.TrimSpace(item) })
	parts = lo.Filter(parts, func(item string, _ int) bool {
		return item != ""
	})
	return parts
}
