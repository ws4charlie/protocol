/*
	Copyright 2017-2018 OneLedger

	Cover over the Tendermint client handling.

	TODO: Make this generic to handle HTTP and local clients
*/
package main

import (
	"os"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	client "github.com/tendermint/abci/client"
	"github.com/tendermint/abci/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var _ *client.Client

// Generic Client interface, allows SetOption
func NewClient() client.Client {

	// TODO: Would like to startup a node, if one isn't running.

	log.Debug("New Client", "address", global.Current.Address, "transport", global.Current.Transport)
	client, err := client.NewClient(global.Current.Address, global.Current.Transport, true)
	if err != nil {
		log.Fatal("Can't start client", "err", err)
	}
	log.Debug("Have Client", "client", client)

	return client
}

func SetOption(key string, value string) {
	client := NewClient()
	options := types.RequestSetOption{
		Key:   key,
		Value: value,
	}

	log.Debug("Setting Option")

	/*
		response := client.SetOptionAsync(options)
	*/

	response, err := client.SetOptionSync(options)
	log.Debug("Have Set Option")

	if err != nil {
		log.Error("SetOption Failed", "err", err, "response", response)
	}
}

var cachedClient *rpcclient.HTTP

// HTTP interface, allows Broadcast?
// TODO: Want to switch client type, based on config or cli args.
func GetClient() *rpcclient.HTTP {
	//cachedClient = rpcclient.NewHTTP("127.0.0.1:46657", "/websocket")
	log.Debug("RPCClient", "address", global.Current.Address)
	cachedClient = rpcclient.NewHTTP(global.Current.Address, "/websocket")
	return cachedClient
}

// Broadcast packet to the chain
func Broadcast(packet []byte) *ctypes.ResultBroadcastTxCommit {
	client := GetClient()

	result, err := client.BroadcastTxCommit(packet)
	if err != nil {
		log.Error("Error", "err", err)
		os.Exit(-1)
	}
	return result
}

// Send a very specific query
func Query(path string, packet []byte) *ctypes.ResultABCIQuery {
	client := GetClient()

	result, err := client.ABCIQuery(path, packet)
	if err != nil {
		log.Error("Error", "err", err)
		os.Exit(-1)
	}
	return result
}