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
			"total_balance":       newGlobalMetric(namespace, "total_balance", "Number of blocks in the longest block chain."),
			"unconfirmed_balance": newGlobalMetric(namespace, "unconfirmed_balance", "Number of active connections to other peers."),
			"confirmed_balance":   newGlobalMetric(namespace, "confirmed_balance", "Proof-of-work difficulty as a multiple of the minimum difficulty"),
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

	ch <- prometheus.MustNewConstMetric(c.metrics["total_balance"],
		prometheus.GaugeValue, float64(stats.Wallet.TotalBallance))
	ch <- prometheus.MustNewConstMetric(c.metrics["unconfirmed_balance"],
		prometheus.GaugeValue, float64(stats.Wallet.UnconfirmedBalance))
	ch <- prometheus.MustNewConstMetric(c.metrics["confirmed_balance"],
		prometheus.GaugeValue, float64(stats.Wallet.ConfirmedBalance))
}
