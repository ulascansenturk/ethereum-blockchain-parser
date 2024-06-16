package ethereum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type EthClient struct {
	url string
}

type RequestBody struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type ClientResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *error          `json:"error"`
}

func NewEthClient(url string) *EthClient {
	return &EthClient{
		url: url,
	}
}
func (rpc *EthClient) Call(method string, params ...interface{}) (json.RawMessage, error) {
	request := RequestBody{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling json: %w", err)
	}

	httpResponse, err := http.Post(rpc.url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer httpResponse.Body.Close()

	data, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var resp ClientResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("rpc error: %w", resp.Error)
	}

	return resp.Result, nil
}

func (rpc *EthClient) callAndUnmarshalResult(method string, target interface{}, params ...interface{}) error {
	result, err := rpc.Call(method, params...)
	if err != nil {
		return fmt.Errorf("error calling method %s: %w", method, err)
	}
	if err := json.Unmarshal(result, target); err != nil {
		return fmt.Errorf("error unmarshaling result: %w", err)
	}
	return nil
}

func (rpc *EthClient) GetBlockNumber() (int, error) {
	var response string
	if err := rpc.callAndUnmarshalResult("eth_blockNumber", &response); err != nil {
		return 0, fmt.Errorf("error getting block number: %w", err)
	}

	blockNumber, err := strconv.ParseInt(response[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing block number: %w", err)
	}

	return int(blockNumber), nil
}

func (rpc *EthClient) GetBlockByNumber(number int) (*Block, error) {
	params := []interface{}{fmt.Sprintf("0x%x", number), true}
	var response Block
	if err := rpc.callAndUnmarshalResult("eth_getBlockByNumber", &response, params...); err != nil {
		return nil, fmt.Errorf("error getting block by number: %w", err)
	}

	return &response, nil
}
