package config

import (
	"cmp"
	"context"
	"fmt"
	"os"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/knadh/koanf/providers/parameterstore/v2"
	"github.com/knadh/koanf/v2"
)

func loadFromAWSParameterStore(k *koanf.Koanf, appName string) error {
	if !ssmEnabled(appName) {
		return nil
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return fmt.Errorf("loading aws config: %w", err)
	}

	// SSM:<prefix>/foo/bar=baz -> foo.bar=baz
	ssmClient := ssm.NewFromConfig(awsCfg)
	ssmCfg := parameterstore.Config[ssm.GetParametersByPathInput]{
		Input: ssm.GetParametersByPathInput{
			Path:           new(ssmPathPrefix(appName)),
			Recursive:      new(true),
			WithDecryption: new(true),
		},
		Delim: "/",
		Callback: func(key, val string) (string, any) {
			if val == "EMPTY" {
				// If the value is "EMPTY", treat it as an explicit empty value.
				return "INVALID_KEY", val
			}

			key = strings.TrimPrefix(key, ssmPathPrefix(appName)+"/")
			return key, val
		},
	}

	psProvider := parameterstore.ProviderWithClient(ssmCfg, ssmClient)
	if err := k.Load(psProvider, nil); err != nil {
		return fmt.Errorf("loading from AWS SSM parameter store: %w", err)
	}

	return nil
}

func ssmEnabled(appName string) bool {
	key := appNamePrefixedKey(appName, "SSM_ENABLED")
	enabled := os.Getenv(key)
	enabled = strings.TrimSpace(enabled)
	enabled = strings.ToLower(enabled)
	return enabled == "1" || enabled == "true" || enabled == "yes"
}

func ssmPathPrefix(appName string) string {
	key := appNamePrefixedKey(appName, "SSM_PATH_PREFIX")
	prefix := os.Getenv(key)
	prefix = cmp.Or(prefix, "/hompy/local")
	prefix = strings.TrimSpace(prefix)
	return prefix
}
