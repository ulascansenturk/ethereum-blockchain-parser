package parser

import (
	"math/big"
	"regexp"
	"strings"
	"ulascansenturk/ethereum-blockchain-parser/pkg/ethereum"
)

func CalculateNextBlock(currentBlock, latestBlock int) int {
	if currentBlock == latestBlock {
		return 0
	}
	if currentBlock == 0 {
		return latestBlock
	}
	return currentBlock + 1
}

func (b *BlockchainScanner) ExtractRelevantTransactions(txs []Transaction) map[string][]Transaction {
	result := make(map[string][]Transaction)
	for _, tx := range txs {
		if exists, _ := b.storage.Has(tx.From); exists {
			result[tx.From] = append(result[tx.From], tx)
		}
		if exists, _ := b.storage.Has(tx.To); exists {
			result[tx.To] = append(result[tx.To], tx)
		}
	}
	return result
}

func ConvertTransactions(txs []ethereum.Transaction) []Transaction {
	transactions := make([]Transaction, len(txs))
	for i, tx := range txs {
		transactions[i] = ParseTransaction(tx)
	}
	return transactions
}

func ParseTransaction(tx ethereum.Transaction) Transaction {
	return Transaction{
		ChainID:     hexToBigInt(tx.ChainID),
		BlockNumber: hexToBigInt(tx.BlockNumber),
		Hash:        tx.Hash,
		Nonce:       hexToBigInt(tx.Nonce),
		From:        tx.From,
		To:          tx.To,
		Value:       hexToBigInt(tx.Value),
		Gas:         hexToBigInt(tx.Gas),
		GasPrice:    hexToBigInt(tx.GasPrice),
		Input:       tx.Input,
	}
}

func hexToBigInt(hexStr string) *big.Int {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	value := new(big.Int)
	value.SetString(hexStr, 16)
	return value
}

func isValidAddress(address string) bool {
	const ethAddressPattern = `^0x[0-9a-fA-F]{40}$`
	match, _ := regexp.MatchString(ethAddressPattern, address)
	return match
}
