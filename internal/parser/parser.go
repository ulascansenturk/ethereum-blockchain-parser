package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"ulascansenturk/ethereum-blockchain-parser/internal/storage"
	"ulascansenturk/ethereum-blockchain-parser/pkg/ethereum"
)

type Service struct {
	storage storage.InMemoryStorage
	*BlockchainScanner
}

type Parser interface {
	GetCurrentBlock() int

	Subscribe(address string) bool

	GetTransactions(address string) []Transaction
}

func New(ctx context.Context, endpoint string, firstBlockToScan int) *Service {
	store := storage.New()
	ethClient := ethereum.NewEthClient(endpoint)
	scanner := NewBlockchainScanner(ctx, store, ethClient, firstBlockToScan)
	return &Service{
		storage:           store,
		BlockchainScanner: scanner,
	}
}

func (s *Service) Subscribe(address string) bool {
	address = strings.ToLower(address)
	if !isValidAddress(address) {
		return false
	}

	if err := s.storage.Put(address, [][]byte{}); err != nil {
		return false
	}
	return true
}

func (s *Service) GetTransactions(address string) []Transaction {
	txs, err := s.storage.Get(strings.ToLower(address))
	if err != nil {
		fmt.Println("error getting transactions: ", err)
		return nil
	}
	transactions := make([]Transaction, len(txs))
	for i, v := range txs {
		var tx Transaction
		if err := json.Unmarshal(v, &tx); err != nil {
			fmt.Println("error unmarshaling transaction: ", err)
			continue
		}
		transactions[i] = tx
	}
	return transactions
}
