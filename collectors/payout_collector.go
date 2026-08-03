package collectors

import (
	"log"

	"github.com/akash329d/storj_exporter/api"
	"github.com/akash329d/storj_exporter/models"
	"github.com/prometheus/client_golang/prometheus"
)

type PayoutCollector struct {
	clients []*api.ApiClient
	metrics map[string]*prometheus.Desc
}

func NewPayoutCollector(clients []*api.ApiClient) *PayoutCollector {
	return &PayoutCollector{
		clients: clients,
		metrics: map[string]*prometheus.Desc{
			"egressBandwidth": prometheus.NewDesc(
				"storj_payout_egress_bandwidth_bytes",
				"Egress bandwidth used by the node for payout calculation",
				[]string{"node_id", "node_name", "period"},
				nil,
			),
			"egressBandwidthPayout": prometheus.NewDesc(
				"storj_payout_egress_bandwidth_cents",
				"Payout for the egress bandwidth used in cents",
				[]string{"node_id", "node_name", "period"},
				nil,
			),
			"egressRepairAudit": prometheus.NewDesc(
				"storj_payout_egress_repair_audit_bytes",
				"Egress bandwidth used for repairs and audits in payout calculation",
				[]string{"node_id", "node_name", "period"},
				nil,
			),
			"egressRepairAuditPayout": prometheus.NewDesc(
				"storj_payout_egress_repair_audit_cents",
				"Payout for the egress bandwidth used for repairs and audits in cents",
				[]string{"node_id", "node_name", "period"},
				nil,
			),
			"diskSpace": prometheus.NewDesc(
				"storj_payout_disk_space_bytes",
				"Disk space used by the node for payout calculation",
				[]string{"node_id", "node_name", "period"},
				nil,
			),
			"diskSpacePayout": prometheus.NewDesc(
				"storj_payout_disk_space_cents",
				"Payout for the disk space used in cents",
				[]string{"node_id", "node_name", "period"},
				nil,
			),
			"heldRate": prometheus.NewDesc(
				"storj_payout_held_rate",
				"Percentage of payout held back by the network",
				[]string{"node_id", "node_name", "period"},
				nil,
			),
			"payout": prometheus.NewDesc(
				"storj_payout_total_cents",
				"Total payout for the node in cents",
				[]string{"node_id", "node_name", "period"},
				nil,
			),
			"held": prometheus.NewDesc(
				"storj_payout_held_cents",
				"Total amount held back by the network in cents",
				[]string{"node_id", "node_name", "period"},
				nil,
			),
			"currentMonthExpectations": prometheus.NewDesc(
				"storj_payout_current_month_expectations_cents",
				"Expected payout for the current month in cents",
				[]string{"node_id", "node_name"},
				nil,
			),
		},
	}
}


func (c *PayoutCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric
	}
}

func (c *PayoutCollector) Collect(ch chan<- prometheus.Metric) {
	for _, client := range c.clients {
		payoutData, err := client.Payout()
		if err != nil {
			log.Printf("Error collecting node payout data: %v", err)
			continue
		}

		c.collectPayoutMetrics(ch, client.NodeID, client.NodeName, payoutData.CurrentMonth, "current")
		c.collectPayoutMetrics(ch, client.NodeID, client.NodeName, payoutData.PreviousMonth, "previous")

		ch <- prometheus.MustNewConstMetric(
			c.metrics["currentMonthExpectations"],
			prometheus.GaugeValue,
			float64(payoutData.CurrentMonthExpectations),
			client.NodeID,
			client.NodeName,
		)
	}
}

func (c *PayoutCollector) collectPayoutMetrics(ch chan<- prometheus.Metric, nodeID string, nodeName string, data models.PayoutData, period string) {
	metrics := map[string]float64{
		"egressBandwidth":         float64(data.EgressBandwidth),
		"egressBandwidthPayout":   float64(data.EgressBandwidthPayout),
		"egressRepairAudit":       float64(data.EgressRepairAudit),
		"egressRepairAuditPayout": float64(data.EgressRepairAuditPayout),
		"diskSpace":               data.DiskSpace,
		"diskSpacePayout":         float64(data.DiskSpacePayout),
		"heldRate":                data.HeldRate * 100, // Convert to percentage
		"payout":                  float64(data.Payout),
		"held":                    float64(data.Held),
	}

	for name, value := range metrics {
		ch <- prometheus.MustNewConstMetric(
			c.metrics[name],
			prometheus.GaugeValue,
			value,
			nodeID,
			nodeName,
			period,
		)
	}
}
