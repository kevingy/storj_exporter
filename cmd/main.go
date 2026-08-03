package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/akash329d/storj_exporter/api"
	"github.com/akash329d/storj_exporter/collectors"
	"github.com/akash329d/storj_exporter/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run() {
	configFile := flag.String("config", "", "Path to a JSON or YAML configuration file")
	flag.Parse()

	nodes := getNodeConfigs(*configFile)
	if len(nodes) == 0 {
		log.Fatal("No Storj node URLs found. Provide a --config file or set STORJ_NODE_<n>_URL environment variables.")
	}

	clients := make([]*api.ApiClient, len(nodes))
	for i, node := range nodes {
		clients[i] = api.NewApiClient(node.URL, node.Name)
	}

	prometheus.MustRegister(collectors.NewNodeCollector(clients))
	prometheus.MustRegister(collectors.NewSatelliteCollector(clients))
	prometheus.MustRegister(collectors.NewPayoutCollector(clients))

	port := 8000 // Default port
	if value, exists := os.LookupEnv("EXPORTER_PORT"); exists {
		if intValue, err := strconv.Atoi(value); err != nil {
			log.Fatalf("Invalid port number in EXPORTER_PORT: %v\n", err)
		} else {
			port = intValue
		}
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting Storj Node Exporter on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// getNodeConfigs returns node configurations from a config file when --config
// is provided, or falls back to STORJ_NODE_<n>_URL environment variables.
func getNodeConfigs(configFile string) []config.NodeConfig {
	if configFile != "" {
		cfg, err := config.Load(configFile)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
		return cfg.Nodes
	}
	return getNodeConfigsFromEnv()
}

// getNodeConfigsFromEnv builds NodeConfig entries from STORJ_NODE_<n>_URL
// environment variables (sequential, starting at 1).
func getNodeConfigsFromEnv() []config.NodeConfig {
	var nodes []config.NodeConfig
	for i := 1; ; i++ {
		rawURL := os.Getenv(fmt.Sprintf("STORJ_NODE_%d_URL", i))
		if rawURL == "" {
			break
		}
		u, err := url.Parse(rawURL)
		if err != nil {
			log.Printf("Error parsing URL for node %d, %s: %v", i, rawURL, err)
			continue
		}
		nodes = append(nodes, config.NodeConfig{URL: u.String()})
	}
	// Apply default names (host:port derived) for env-based nodes.
	cfg := &config.Config{Nodes: nodes}
	config.ApplyDefaults(cfg)
	return cfg.Nodes
}