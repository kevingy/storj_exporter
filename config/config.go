package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// NodeConfig represents a single Storj node entry in the configuration file.
type NodeConfig struct {
	URL  string `json:"url"  yaml:"url"`
	Name string `json:"name" yaml:"name"`
}

// Config represents the top-level configuration file structure.
type Config struct {
	Nodes []NodeConfig `json:"nodes" yaml:"nodes"`
}

// Load reads a configuration file from path, auto-detects its format based on
// the file extension (.json, .yaml, .yml), and returns the parsed Config.
// If the extension is unrecognised, or the file content is invalid for the
// detected format, a descriptive error is returned.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	ext := strings.ToLower(filepath.Ext(path))

	var cfg Config
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parsing JSON config file %q: %w", path, err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parsing YAML config file %q: %w", path, err)
		}
	default:
		return nil, fmt.Errorf("cannot determine config file format for %q: use a .json, .yaml, or .yml extension", path)
	}

	if err := validate(&cfg, path); err != nil {
		return nil, err
	}

	applyDefaults(&cfg)
	return &cfg, nil
}

// validate checks that each node entry has a non-empty URL.
func validate(cfg *Config, path string) error {
	for i, node := range cfg.Nodes {
		if strings.TrimSpace(node.URL) == "" {
			return fmt.Errorf("config file %q: node[%d] is missing required field \"url\"", path, i)
		}
	}
	return nil
}

// applyDefaults fills in the Name field for nodes where it was omitted,
// deriving the name from the URL host and port (e.g. "192.168.1.10:14002").
func applyDefaults(cfg *Config) {
	for i := range cfg.Nodes {
		if strings.TrimSpace(cfg.Nodes[i].Name) != "" {
			continue
		}
		cfg.Nodes[i].Name = deriveNameFromURL(cfg.Nodes[i].URL)
	}
}

// deriveNameFromURL returns host:port (or just host when no port is present)
// extracted from rawURL.  If parsing fails the raw URL string is returned.
func deriveNameFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	host := u.Host // already "host:port" or just "host"
	if host == "" {
		return rawURL
	}
	return host
}
