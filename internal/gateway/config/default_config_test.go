package config_test

import (
	"os"
	"testing"

	"github.com/knadh/koanf/parsers/yaml"
)

// The image ships configs/gateway/config.default.yaml and the loader merges it under every
// deployment's config, so a key pair there would be trusted everywhere.
func TestDefaultConfigShipsNoJWTKeyPair(t *testing.T) {
	raw, err := os.ReadFile("../../../configs/gateway/config.default.yaml")
	if err != nil {
		t.Fatal(err)
	}
	doc, err := yaml.Parser().Unmarshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	auth, _ := doc["auth"].(map[string]any)
	jwt, _ := auth["jwt"].(map[string]any)
	if pairs, _ := jwt["key-pairs"].(map[string]any); len(pairs) != 0 {
		t.Fatalf("config.default.yaml has JWT key pairs %v; keys must come from each deployment", keys(pairs))
	}
}

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
