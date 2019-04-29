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

// Stats represents node metrics.
type Stats struct {
	Wallet StubWallet
	Node   StubNode
}

type StubWallet struct {
	TotalBallance      int64
	ConfirmedBalance   int64
	UnconfirmedBalance int64
}

type StubNode struct {
	Peers            uint32
	PendingChannels  uint32
	ActiveChannels   uint32
	InactiveChannels uint32
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
func (client *LightningClient) GetStats() (*Stats, error) {

	var stats Stats

	ctxb := context.Background()

	reqWallet := &lnrpc.WalletBalanceRequest{}
	wallet, err := client.rpcclient.WalletBalance(ctxb, reqWallet)
	if err != nil {
		log.Fatal(err)
	}

	// Wallet
	totalBalance := wallet.TotalBalance
	stats.Wallet.TotalBallance = totalBalance

	unconfirmedBalance := wallet.UnconfirmedBalance
	stats.Wallet.UnconfirmedBalance = unconfirmedBalance

	confirmedBalance := wallet.ConfirmedBalance
	stats.Wallet.ConfirmedBalance = confirmedBalance

	// Info
	reqInfo := &lnrpc.GetInfoRequest{}
	info, err := client.rpcclient.GetInfo(ctxb, reqInfo)
	if err != nil {
		log.Fatal(err)
	}
	peers := info.NumPeers
	stats.Node.Peers = peers

	numInactiveChannels := info.NumInactiveChannels
	stats.Node.InactiveChannels = numInactiveChannels

	numActiveChannels := info.NumActiveChannels
	stats.Node.ActiveChannels = numActiveChannels

	numPendingChannels := info.NumPendingChannels
	stats.Node.PendingChannels = numPendingChannels

	return &stats, nil
}
