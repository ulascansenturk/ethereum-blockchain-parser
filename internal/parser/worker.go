package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"ulascansenturk/ethereum-blockchain-parser/internal/storage"
	"ulascansenturk/ethereum-blockchain-parser/pkg/ethereum"
)

type BlockchainScanner struct {
	ctx              context.Context
	storage          storage.InMemoryStorage
	ethClient        *ethereum.EthClient
	lastScannedBlock int
}

type Worker interface {
	Start() (int, error)
	GetCurrentBlock() int
}

func NewBlockchainScanner(ctx context.Context, storage storage.InMemoryStorage, ethClient *ethereum.EthClient, firstBlock int) *BlockchainScanner {
	return &BlockchainScanner{
		ctx:              ctx,
		storage:          storage,
		ethClient:        ethClient,
		lastScannedBlock: firstBlock,
	}
}

func (b *BlockchainScanner) GetCurrentBlock() int {
	return b.lastScannedBlock
}

func (b *BlockchainScanner) Start() (int, error) {
	headBlock, err := b.ethClient.GetBlockNumber()
	if err != nil {
		log.Printf("Error querying head block number: %v", err)
		return 0, err
	}

	nextBlock := CalculateNextBlock(b.lastScannedBlock, headBlock)
	if nextBlock == 0 {
		return 0, nil
	}

	txs, err := b.processBlock(nextBlock)
	if err != nil {
		log.Printf("Error processing block %d: %v", nextBlock, err)
		return 0, err
	}

	if len(txs) > 0 {
		if err := b.storeTransactions(txs); err != nil {
			log.Printf("Error storing transactions: %v", err)
			return 0, err
		}
	}

	b.lastScannedBlock = nextBlock
	return b.lastScannedBlock, nil
}

func (b *BlockchainScanner) StartScan(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-b.ctx.Done():
			return
		case <-ticker.C:
			if _, err := b.scanBlocks(); err != nil {
				log.Printf("Error scanning blocks: %v", err)
			}
			log.Printf("Last scanned block %d", b.GetCurrentBlock())
		}
	}
}

func (b *BlockchainScanner) scanBlocks() (int, error) {
	for {
		scannedBlock, err := b.Start()
		if scannedBlock == 0 || err != nil {
			return scannedBlock, err
		}
	}
}

func (b *BlockchainScanner) processBlock(blockNumber int) (map[string][]Transaction, error) {
	block, err := b.ethClient.GetBlockByNumber(blockNumber)
	if err != nil {
		log.Printf("Error querying block %d: %v", blockNumber, err)
		return nil, err
	}

	transactions := b.ExtractRelevantTransactions(ConvertTransactions(block.Transactions))
	if len(transactions) == 0 {
		return nil, nil
	}

	return transactions, nil
}

func (b *BlockchainScanner) storeTransactions(transactions map[string][]Transaction) error {
	for address, txs := range transactions {
		data, err := encodeTransactions(txs)
		if err != nil {
			return fmt.Errorf("error encoding transactions for address %s: %w", address, err)
		}
		if err := b.storage.Put(address, data); err != nil {
			return fmt.Errorf("error saving transactions for address %s: %w", address, err)
		}
	}
	return nil
}

func encodeTransactions(batch []Transaction) ([][]byte, error) {
	var txs [][]byte
	for _, tx := range batch {
		data, err := json.Marshal(tx)
		if err != nil {
			return nil, fmt.Errorf("error marshaling transaction: %w", err)
		}
		txs = append(txs, data)
	}
	return txs, nil
}
