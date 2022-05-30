package gocore

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/core-coin/go-core/common/hexutil"
)

type RPCClient struct {
	Url    string
	client *http.Client

	offline bool
}

type JSONRpcResp struct {
	Id     *json.RawMessage       `json:"id"`
	Result *json.RawMessage       `json:"result"`
	Error  map[string]interface{} `json:"error"`
}

func NewRPCClient(url string, timeout time.Duration) *RPCClient {
	rpcClient := &RPCClient{Url: url}
	rpcClient.client = &http.Client{
		Timeout: timeout,
	}
	if url == "" {
		rpcClient.offline = true
	}
	return rpcClient
}

func (r *RPCClient) doPost(url string, method string, params interface{}) (*JSONRpcResp, error) {
	if r.offline {
		return nil, errors.New("You cannot do it without connection to gocore RPC API, try to add flag --gocore")
	}
	jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	data, err := json.Marshal(jsonReq)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp *JSONRpcResp
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		return nil, err
	}
	if rpcResp.Error != nil {
		return nil, errors.New(rpcResp.Error["message"].(string))
	}
	return rpcResp, err
}

func (r *RPCClient) SendRawTransaction(data string) (string, error) {
	params := []string{data}
	rpcResp, err := r.doPost(r.Url, "xcb_sendRawTransaction", params)
	var reply string
	if err != nil {
		return reply, err
	}
	err = json.Unmarshal(*rpcResp.Result, &reply)
	if err != nil {
		return reply, err
	}
	return reply, err
}

func (r *RPCClient) GetAccountNonce(account, status string) (uint64, error) {
	params := []string{account, status}

	rpcResp, err := r.doPost(r.Url, "xcb_getTransactionCount", params)
	if err != nil {
		return 0, err
	}
	var reply string
	err = json.Unmarshal(*rpcResp.Result, &reply)
	resInt, err := strconv.ParseInt(reply[2:], 16, 64)
	if err != nil {
		return 0, err
	}
	return uint64(resInt), err
}

func (r *RPCClient) EstimateEnergyPrice() (int64, error) {
	rpcResp, err := r.doPost(r.Url, "xcb_energyPrice", nil)
	if err != nil {
		return 0, err
	}
	var reply string
	err = json.Unmarshal(*rpcResp.Result, &reply)
	if err != nil {
		return 0, err
	}
	energy, err := hexutil.DecodeUint64(reply)
	if err != nil {
		return 0, err
	}
	return int64(energy), nil
}
