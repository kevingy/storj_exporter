package config

import (
	"os"
	"path/filepath"
	"testing"
)

// writeFile is a helper that creates a temporary file with the given content
// and extension, returning its path.
func writeFile(t *testing.T, ext, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "storj_config_*"+ext)
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadJSON(t *testing.T) {
	content := `{
		"nodes": [
			{"url": "http://192.168.1.10:14002", "name": "my-node"},
			{"url": "http://192.168.1.11:14002"}
		]
	}`
	path := writeFile(t, ".json", content)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(cfg.Nodes))
	}
	if cfg.Nodes[0].Name != "my-node" {
		t.Errorf("node[0] name: got %q, want %q", cfg.Nodes[0].Name, "my-node")
	}
	// node[1] name should be derived from URL
	if cfg.Nodes[1].Name != "192.168.1.11:14002" {
		t.Errorf("node[1] derived name: got %q, want %q", cfg.Nodes[1].Name, "192.168.1.11:14002")
	}
}

func TestLoadYAML(t *testing.T) {
	content := `
nodes:
  - url: "http://example.com:14002"
    name: "example-node"
  - url: "http://192.168.1.20:14002"
`
	path := writeFile(t, ".yaml", content)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(cfg.Nodes))
	}
	if cfg.Nodes[0].Name != "example-node" {
		t.Errorf("node[0] name: got %q, want %q", cfg.Nodes[0].Name, "example-node")
	}
	if cfg.Nodes[1].Name != "192.168.1.20:14002" {
		t.Errorf("node[1] derived name: got %q, want %q", cfg.Nodes[1].Name, "192.168.1.20:14002")
	}
}

func TestLoadYMLExtension(t *testing.T) {
	content := `
nodes:
  - url: "http://node.example.com:14002"
`
	path := writeFile(t, ".yml", content)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(cfg.Nodes))
	}
	if cfg.Nodes[0].Name != "node.example.com:14002" {
		t.Errorf("derived name: got %q, want %q", cfg.Nodes[0].Name, "node.example.com:14002")
	}
}

func TestNameOmittedFallback(t *testing.T) {
	content := `{"nodes": [{"url": "http://myhost:9999"}]}`
	path := writeFile(t, ".json", content)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Nodes[0].Name != "myhost:9999" {
		t.Errorf("expected derived name %q, got %q", "myhost:9999", cfg.Nodes[0].Name)
	}
}

func TestInvalidJSON(t *testing.T) {
	path := writeFile(t, ".json", `{invalid json`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestInvalidYAML(t *testing.T) {
	path := writeFile(t, ".yaml", "nodes:\n  - url: [\nbad")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestUndeterminedFormat(t *testing.T) {
	path := writeFile(t, ".toml", `url = "http://localhost:14002"`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for undetermined format, got nil")
	}
}

func TestMissingURL(t *testing.T) {
	content := `{"nodes": [{"name": "no-url-node"}]}`
	path := writeFile(t, ".json", content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing URL, got nil")
	}
}

func TestFileNotFound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
