package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"ulascansenturk/ethereum-blockchain-parser/internal/parser"
)

const (
	EthRpcURL       = "https://ethereum-rpc.publicnode.com"
	ServerAddr      = ":8080"
	ScanInterval    = 12 * time.Second
	ContentTypeJSON = "application/json"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := parser.New(ctx, EthRpcURL, 0)
	go service.StartScan(ScanInterval)

	http.HandleFunc("/subscribe", handleSubscribe(service))
	http.HandleFunc("/transactions", handleGetTransactions(service))
	http.HandleFunc("/currentBlock", handleGetCurrentBlock(service))

	log.Printf("Starting server on %s...", ServerAddr)
	if err := http.ListenAndServe(ServerAddr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleSubscribe(service *parser.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}

		subscribed := service.Subscribe(address)
		if !subscribed {
			http.Error(w, "Failed to subscribe address", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Subscribed successfully"))
	}
}

func handleGetTransactions(service *parser.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}

		transactions := service.GetTransactions(address)
		if transactions == nil {
			http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(transactions)
		if err != nil {
			http.Error(w, "Failed to marshal transactions", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", ContentTypeJSON)
		w.Write(response)
	}
}

func handleGetCurrentBlock(service *parser.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		currentBlock := service.GetCurrentBlock()
		response := map[string]int{"current_block": currentBlock}

		respJson, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", ContentTypeJSON)
		w.Write(respJson)
	}
}
