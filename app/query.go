/*

   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package app

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/tendermint/tendermint/abci/types"
)

func (app *App) Query(req RequestQuery) ResponseQuery {
	app.logger.Debug("ABCI: Query", "req", req, "path", req.Path, "data", req.Data)

	routerReq := NewRequestFromData(req.Path, req.Data)
	routerReq.Ctx = app.Context

	resp := &Response{}
	app.r.Handle(*routerReq, resp)

	// TODO proper response error handling
	result := ResponseQuery{
		Code:  types.CodeTypeOK,
		Index: 0, // TODO: What is this for?

		Log:  "Log Information",
		Info: "Info Information",

		Key:   action.Message("result"),
		Value: resp.Data,

		Proof:  nil,
		Height: int64(app.Context.balances.Version),
	}

	app.logger.Debug("ABCI: Query Result", "result", result)
	return result
}

func NewABCIRouter() Router {
	r := NewRouter("abci")
	r.AddHandler("/query/balance", GetBalance)
}

/*
		Handlers start here
 */
func GetBalance(req Request, resp *Response) {
	req.Parse()
	key := req.GetBytes("key")
	if len(key) == 0 {
		resp.Error("required parameter key")
	}


}
