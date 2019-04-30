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

type WalletStats struct {
	TotalBallance      int64
	ConfirmedBalance   int64
	UnconfirmedBalance int64
}

type NodeStats struct {
	Peers            uint32
	PendingChannels  uint32
	ActiveChannels   uint32
	InactiveChannels uint32
	BlockHeight      uint32
	SyncedToChain    uint8
}

type PendingChannelsStats struct {
	TotalLimboBalance           int64
	PendingOpenChannels         int
	PendingClosingChannels      int
	PendingForceClosingChannels int
	WaitingCloseChannels        int
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
func (client *LightningClient) GetStats() (*lnrpc.GetInfoResponse, error) {
	ctxb := context.Background()

	// Pending Channels
	req := &lnrpc.GetInfoRequest{}
	info, err := client.rpcclient.GetInfo(ctxb, req)
	if err != nil {
		log.Fatal(err)
	}

	return info, err
}

// GetWalletStats get wallet balances
func (client *LightningClient) GetWalletStats() (*WalletStats, error) {
	var stats WalletStats

	ctxb := context.Background()

	req := &lnrpc.WalletBalanceRequest{}
	wallet, err := client.rpcclient.WalletBalance(ctxb, req)
	if err != nil {
		log.Fatal(err)
	}

	stats.TotalBallance = wallet.TotalBalance
	stats.UnconfirmedBalance = wallet.UnconfirmedBalance
	stats.ConfirmedBalance = wallet.ConfirmedBalance

	return &stats, nil
}

// GetInfoStats gets general node info
func (client *LightningClient) GetInfoStats() (*NodeStats, error) {
	var stats NodeStats

	ctxb := context.Background()

	req := &lnrpc.GetInfoRequest{}
	info, err := client.rpcclient.GetInfo(ctxb, req)
	if err != nil {
		log.Fatal(err)
	}
	stats.Peers = info.NumPeers
	stats.InactiveChannels = info.NumInactiveChannels
	stats.ActiveChannels = info.NumActiveChannels
	stats.PendingChannels = info.NumPendingChannels
	stats.BlockHeight = info.BlockHeight
	stats.SyncedToChain = boolToInt(info.SyncedToChain)

	return &stats, nil
}

// GetPendingChannelsStats get pending channels status
func (client *LightningClient) GetPendingChannelsStats() (*PendingChannelsStats, error) {
	var stats PendingChannelsStats

	ctxb := context.Background()

	req := &lnrpc.PendingChannelsRequest{}
	info, err := client.rpcclient.PendingChannels(ctxb, req)
	if err != nil {
		log.Fatal(err)
	}

	stats.TotalLimboBalance = info.TotalLimboBalance
	stats.PendingOpenChannels = len(info.PendingOpenChannels)
	stats.PendingClosingChannels = len(info.PendingClosingChannels)
	stats.PendingForceClosingChannels = len(info.PendingForceClosingChannels)
	stats.WaitingCloseChannels = len(info.WaitingCloseChannels)

	return &stats, nil
}

func boolToInt(arg bool) uint8 {
	if arg {
		return 1
	}
	return 0
}
