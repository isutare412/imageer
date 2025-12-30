package config

import (
	"github.com/isutare412/imageer/internal/gateway/crypt"
	"github.com/isutare412/imageer/internal/gateway/jwt"
	"github.com/isutare412/imageer/internal/gateway/oidc"
	"github.com/isutare412/imageer/internal/gateway/postgres"
	"github.com/isutare412/imageer/internal/gateway/service/auth"
	"github.com/isutare412/imageer/internal/gateway/web"
	"github.com/isutare412/imageer/pkg/log"
)

func (c *Config) ToLogConfig() log.Config {
	return log.Config(c.Log)
}

func (c *Config) ToWebConfig() web.Config {
	return web.Config{
		Port:            c.Web.Port,
		ShowBanner:      c.Web.ShowBanner,
		ShowOpenAPIDocs: c.Web.ShowOpenAPIDocs,
		APIKeyHeader:    c.Auth.ServiceAccount.APIKeyHeader,
		UserCookieName:  c.Auth.Cookies.User.Name,
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
