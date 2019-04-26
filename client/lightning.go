package client

import (
	"context"
	"fmt"
	"log"

	"github.com/lightningnetwork/lnd/lnrpc"
)

// LightningClient allows you to fetch Lightning node metrics from rpc.
type LightningClient struct {
	rpcclient lnrpc.LightningClient
}

// NodeStats represents node metrics.
type NodeStats struct {
	Wallet StubWallet
	// ConnectionCount int64
	// Difficulty      float64
}

type StubWallet struct {
	TotalBallance      int64
	ConfirmedBalance   int64
	UnconfirmedBalance int64
}

// NewLightningClient creates an LightningClient.
func NewLightningClient(rpcclient lnrpc.LightningClient) (*LightningClient, error) {

	client := &LightningClient{
		rpcclient: rpcclient,
	}

	if _, err := client.GetStats(); err != nil {
		return nil, fmt.Errorf("Failed to create LightningClient: %v", err)
	}

	return client, nil
}

// GetStats fetches the node metrics.
func (client *LightningClient) GetStats() (*NodeStats, error) {

	var stats NodeStats

	ctxb := context.Background()

	req := &lnrpc.WalletBalanceRequest{}
	resp, err := client.rpcclient.WalletBalance(ctxb, req)
	if err != nil {
		log.Fatal(err)
	}
	totalBalance := resp.TotalBalance
	stats.Wallet.TotalBallance = totalBalance

	unconfirmedBalance := resp.UnconfirmedBalance
	stats.Wallet.UnconfirmedBalance = unconfirmedBalance

	confirmedBalance := resp.ConfirmedBalance
	stats.Wallet.ConfirmedBalance = confirmedBalance

	return &stats, nil
}
