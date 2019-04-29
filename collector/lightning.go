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
			"wallet_balance_satoshis": newGlobalMetric(namespace, "wallet_balance_satoshis", "The wallet balance.", []string{"status"}),
			"peers":                   newGlobalMetric(namespace, "peers", "Number of currently connected peers.", []string{}),
			"channels":                newGlobalMetric(namespace, "channels", "Number of channels", []string{"status"}),
			"block_height":            newGlobalMetric(namespace, "block_height", "The node’s current view of the height of the best block", []string{}),
			"synced_to_chain":         newGlobalMetric(namespace, "synced_to_chain", "The node’s current view of the height of the best block", []string{}),
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

	ch <- prometheus.MustNewConstMetric(c.metrics["wallet_balance_satoshis"],
		prometheus.GaugeValue, float64(stats.Wallet.UnconfirmedBalance), "unconfirmed")
	ch <- prometheus.MustNewConstMetric(c.metrics["wallet_balance_satoshis"],
		prometheus.GaugeValue, float64(stats.Wallet.ConfirmedBalance), "confirmed")

	ch <- prometheus.MustNewConstMetric(c.metrics["peers"],
		prometheus.GaugeValue, float64(stats.Node.Peers))
	ch <- prometheus.MustNewConstMetric(c.metrics["channels"],
		prometheus.GaugeValue, float64(stats.Node.ActiveChannels), "active")
	ch <- prometheus.MustNewConstMetric(c.metrics["channels"],
		prometheus.GaugeValue, float64(stats.Node.PendingChannels), "pending")
	ch <- prometheus.MustNewConstMetric(c.metrics["channels"],
		prometheus.GaugeValue, float64(stats.Node.InactiveChannels), "inactive")
	ch <- prometheus.MustNewConstMetric(c.metrics["block_height"],
		prometheus.GaugeValue, float64(stats.Node.BlockHeight))
	ch <- prometheus.MustNewConstMetric(c.metrics["synced_to_chain"],
		prometheus.GaugeValue, float64(stats.Node.SyncedToChain))

}
