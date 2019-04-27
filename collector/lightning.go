package collector

import (
	"log"
	"sync"

	"github.com/platanus/lightning-prometheus-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

// LightningCollector collects node metrics. It implements prometheus.Collector interface.
type LightningCollector struct {
	lightningClient *client.LightningClient
	metrics         map[string]*prometheus.Desc
	mutex           sync.Mutex
}

// NewLightningCollector creates an LightningCollector.
func NewLightningCollector(lightningClient *client.LightningClient, namespace string) *LightningCollector {
	return &LightningCollector{
		lightningClient: lightningClient,
		metrics: map[string]*prometheus.Desc{
			"wallet_balance_satoshis_total": newGlobalMetric(namespace, "wallet_balance_satoshis_total", "The wallet balance.", []string{"status"}),
		},
	}
}

// Describe sends the super-set of all possible descriptors of node metrics
// to the provided channel.
func (c *LightningCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

// Collect fetches metrics from the node and sends them to the provided channel.
func (c *LightningCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // To protect metrics from concurrent collects
	defer c.mutex.Unlock()

	stats, err := c.lightningClient.GetStats()
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.metrics["wallet_balance_satoshis_total"],
		prometheus.GaugeValue, float64(stats.Wallet.UnconfirmedBalance), "unconfirmed")
	ch <- prometheus.MustNewConstMetric(c.metrics["wallet_balance_satoshis_total"],
		prometheus.GaugeValue, float64(stats.Wallet.ConfirmedBalance), "confirmed")
}
