/*

 */

package btc

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type JobBTCBroadcast struct {
	Type string

	TrackerName string

	JobID string

	BroadcastSuccessful bool

	Done bool

	RetryCount int8
}

func (j *JobBTCBroadcast) DoMyJob(ctxI interface{}) {

	ctx, _ := ctxI.(*action.JobsContext)

	tracker, err := ctx.Trackers.Get(j.TrackerName)
	if err != nil {
		j.RetryCount += 1
		return
	}

	if !j.BroadcastSuccessful {

		lockTx := wire.NewMsgTx(wire.TxVersion)
		err = lockTx.Deserialize(bytes.NewReader(tracker.ProcessUnsignedTx))
		if err != nil {
			j.RetryCount += 1
			return
		}

		type sign []byte
		btcSignatures := tracker.Multisig.GetSignatures()
		signatures := make([]sign, len(btcSignatures))
		for i := range btcSignatures {
			index := btcSignatures[i].Index
			signatures[index] = btcSignatures[i].Sign
		}

		builder := txscript.NewScriptBuilder().AddOp(txscript.OP_FALSE)
		for i := range signatures {
			builder.AddData(signatures[i])
			if i == tracker.Multisig.M {
				break
			}
		}

		lockScript, err := ctx.LockScripts.GetLockScript(tracker.CurrentLockScriptAddress)
		if err != nil {
			j.RetryCount += 1
			return
		}

		builder.AddFullData(lockScript)
		sigScript, err := builder.Script()

		cd := bitcoin.NewChainDriver(ctx.BlockCypherToken)
		lockTx = cd.AddLockSignature(tracker.ProcessUnsignedTx, sigScript)

		buf := bytes.NewBuffer([]byte{})
		lockTx.Serialize(buf)

		// TODO load from config
		connCfg := &rpcclient.ConnConfig{
			Host:         ctx.BTCNodeAddress + ":" + ctx.BTCRPCPort,
			User:         ctx.BTCRPCUsername,
			Pass:         ctx.BTCRPCPassword,
			HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
			DisableTLS:   true, // Bitcoin core does not provide TLS by default
		}

		clt, err := rpcclient.New(connCfg, nil)
		if err != nil {
			j.RetryCount += 1
			return
		}

		var txBytes []byte
		buf = bytes.NewBuffer(txBytes)
		lockTx.Serialize(buf)
		txBytes = buf.Bytes()

		fmt.Println("--------------------------------------------------------------")
		fmt.Println(hex.EncodeToString(txBytes))
		fmt.Println(hex.EncodeToString(lockTx.TxIn[0].SignatureScript))

		hash, err := cd.BroadcastTx(lockTx, clt)
		if err == nil {
			fmt.Println("btc tx hash", hash)
			j.BroadcastSuccessful = true
			return
		} else {

			fmt.Println("broadcast failed, but going forward")
			j.BroadcastSuccessful = true
			j.RetryCount += 1
			return
		}

	} else {

		cd := bitcoin.NewChainDriver(ctx.BlockCypherToken)

		chain := "test3"
		switch ctx.BTCChainnet {
		case "testnet3":
			chain = "test3"
		case "testnet":
			chain = "test"
		case "mainnet":
			chain = "main"
		}

		tempHash, _ := chainhash.NewHashFromStr("860a32ef84ed54df86d207112d1f8d3d5ad28751b25cc7e2107ef55cccbc7586")

		ok, _ := cd.CheckFinality(tempHash, ctx.BlockCypherToken, chain)
		//	ok, _ := cd.CheckFinality(tracker.ProcessTxId, ctx.BlockCypherToken, chain)
		if !ok {
			j.RetryCount += 1
			// return
		}

		data := [4]byte{}
		_, err = io.ReadFull(rand.Reader, data[:])
		if err != nil {
			j.RetryCount += 1
			return
		}

		reportFinalityMint := ReportFinalityMint{
			TrackerName:      j.TrackerName,
			OwnerAddress:     tracker.ProcessOwner,
			ValidatorAddress: ctx.ValidatorAddress,
			RandomBytes:      data[:],
		}

		fmt.Println(0)

		txData, err := reportFinalityMint.Marshal()
		if err != nil {
			// retry later
			j.RetryCount += 1
			return
		}

		fmt.Println(1)
		tx := action.RawTx{
			Type: action.BTC_REPORT_FINALITY_MINT,
			Data: txData,
			Fee:  action.Fee{},
			Memo: j.JobID,
		}

		req := action.InternalBroadcastRequest{
			RawTx: tx,
		}
		rep := action.BroadcastReply{}

		err = ctx.Service.InternalBroadcast(req, &rep)
		if err != nil {

			fmt.Println("internal broadcast error ", err)
			// retry later
			j.RetryCount += 1
			return
		}

	}
}

func (j *JobBTCBroadcast) IsMyJobDone(ctxI interface{}) bool {
	ctx, _ := ctxI.(*action.JobsContext)

	if j.RetryCount > 20 {
		return true
	}

	tracker, err := ctx.Trackers.Get(j.TrackerName)
	if err != nil {
		return false
	}

	if tracker.IsAvailable() {
		return true
	}

	for _, fv := range tracker.FinalityVotes {
		if bytes.Equal(fv, ctx.ValidatorAddress) {

			return true
		}
	}

	return false
}

func (j *JobBTCBroadcast) IsSufficient(ctxI interface{}) bool {
	return j.IsMyJobDone(ctxI)
}

func (j *JobBTCBroadcast) DoFinalize() {
	j.Done = true
}

/*
	simple getters
*/
func (j *JobBTCBroadcast) GetType() string {
	return JobTypeBTCBroadcast
}

func (j *JobBTCBroadcast) GetJobID() string {
	return j.JobID
}

func (j *JobBTCBroadcast) IsDone() bool {
	return j.Done
}
