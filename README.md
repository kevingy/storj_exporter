# Storj Exporter for Prometheus

A Go-based exporter to pull information from Storj node APIs and export it for Prometheus monitoring. Supports monitoring multiple nodes, with each metric labeled by `node_id` and `node_name`.

## Docker Usage

### Docker Hub Repository

Pull the Docker image from [Docker Hub](https://hub.docker.com/r/akash329d/storj_exporter).

### Example: Running the Docker Container

```sh
docker run -d \
  --name=StorjExporter \
  -e STORJ_NODE_1_URL=http://192.168.1.10:14002/ \
  -p 8000:8000 \
  akash329d/storj_exporter
```

### Docker Parameters

| Parameter           | Description                                                                         | Default Value |
|---------------------|-------------------------------------------------------------------------------------|---------------|
| `EXPORTER_PORT`     | Port for the metrics server.                                                        | 8000          |
| `STORJ_NODE_%d_URL` | URL of a Storj node (replace `%d` with a sequential number starting at 1)          | N/A           |

## Configuration File (Multi-Node)

Instead of environment variables, you can use a JSON or YAML configuration file to define multiple nodes with optional friendly names. Pass the file path via the `--config` flag.

### JSON Example (`nodes.json`)

```json
{
  "nodes": [
    { "url": "http://192.168.1.10:14002", "name": "node-1" },
    { "url": "http://192.168.1.11:14002", "name": "node-2" },
    { "url": "http://192.168.1.12:14002" }
  ]
}
```

### YAML Example (`nodes.yaml`)

```yaml
nodes:
  - url: "http://192.168.1.10:14002"
    name: "node-1"
  - url: "http://192.168.1.11:14002"
    name: "node-2"
  - url: "http://192.168.1.12:14002"
```

> When `name` is omitted, it defaults to `host:port` derived from the URL.

### Running with a Config File

```sh
./storj_exporter --config nodes.yaml
```

Or with Docker:

```sh
docker run -d \
  --name=StorjExporter \
  -v /path/to/nodes.yaml:/etc/storj/nodes.yaml \
  -p 8000:8000 \
  akash329d/storj_exporter --config /etc/storj/nodes.yaml
```

## Accessing Metrics

Access the metrics at:
```arduino
http://<host>:8000/metrics
```

Replace <host> with the IP address or hostname of the machine running the Docker container.

## Prometheus Configuration
Add the following job to your prometheus.yml:
```yaml
scrape_configs:
  - job_name: 'storj'
    static_configs:
      - targets: ['http://<host>:8000']
```

## Grafana Dashboard (Metrics Redacted)
![Screenshot](https://github.com/akash329d/storj_exporter/blob/main/grafana_redacted.png?raw=true)

A pre-configured Grafana dashboard is included in this repository as `grafana_dashboard.json`. To use it:

1. Import the JSON file into your Grafana instance.
2. Configure the dashboard to use your Prometheus data source.