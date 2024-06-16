package parser

import "math/big"

type Transaction struct {
	ChainID     *big.Int `json:"chainId"`
	BlockNumber *big.Int `json:"blockNumber"`
	Hash        string   `json:"hash"`
	Nonce       *big.Int `json:"nonce"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	Value       *big.Int `json:"value"`
	Gas         *big.Int `json:"gas"`
	GasPrice    *big.Int `json:"gasPrice"`
	Input       string   `json:"input"`
}
