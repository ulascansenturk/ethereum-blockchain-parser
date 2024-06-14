package main

import (
	"fmt"
	"ulascansenturk/ethereum-blockchain-parser/pkg/ethereum"
)

func main() {
	client := ethereum.NewEthClient("https://cloudflare-eth.com")
	bn, _ := client.GetBlockNumber()
	fmt.Println(client.GetBlockByNumber(bn))

}
