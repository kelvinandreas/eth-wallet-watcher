package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/helper"
)

type EtherscanService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewEtherscanService(apiKey, baseURL string) *EtherscanService {
	return &EtherscanService{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

type EtherscanTx struct {
	Hash            string    `json:"hash"`
	From            string    `json:"from"`
	To              string    `json:"to"`
	Value           string    `json:"value"`
	BlockNumber     string    `json:"block_number"`
	TimeStamp       time.Time `json:"time_stamp"`
	Gas             string    `json:"gas"`
	GasPrice        string    `json:"gas_price"`
	IsError         string    `json:"is_error"`
	TxReceiptStatus string    `json:"txreceipt_status"`
}

type etherscanResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Result  []rawTx `json:"result"`
}

type rawTx struct {
	BlockNumber     string `json:"blockNumber"`
	TimeStamp       string `json:"timeStamp"`
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	Gas             string `json:"gas"`
	GasPrice        string `json:"gasPrice"`
	IsError         string `json:"isError"`
	TxReceiptStatus string `json:"txreceipt_status"`
}

type FetchTransactionsResult struct {
	Transactions []EtherscanTx
	LastBlock    string
}

func (s *EtherscanService) FetchTransactions(address string, lastBlock string) (*FetchTransactionsResult, error) {
	if lastBlock == "" {
		lastBlock = "0"
	}

	params := map[string]string{
		"chainid":    "1",
		"module":     "account",
		"action":     "txlist",
		"address":    address,
		"startblock": lastBlock,
		"endblock":   "99999999",
		"sort":       "asc",
		"apikey":     s.apiKey,
	}

	requestURL, err := helper.BuildURLWithQuery(s.baseURL, params)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("etherscan http error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var er etherscanResponse
	if err := json.Unmarshal(body, &er); err != nil {
		return nil, err
	}

	if er.Status != "1" {
		if er.Message == "No transactions found" || er.Message == "No records found" {
			return &FetchTransactionsResult{
				Transactions: []EtherscanTx{},
				LastBlock:    lastBlock,
			}, nil
		}
		return nil, fmt.Errorf("etherscan error: %s", er.Message)
	}

	out := make([]EtherscanTx, 0, len(er.Result))

	for _, r := range er.Result {
		tsInt, err := strconv.ParseInt(r.TimeStamp, 10, 64)
		if err != nil {
			continue
		}
		t := time.Unix(tsInt, 0)

		tx := EtherscanTx{
			Hash:            r.Hash,
			From:            r.From,
			To:              r.To,
			Value:           r.Value,
			BlockNumber:     r.BlockNumber,
			TimeStamp:       t,
			Gas:             r.Gas,
			GasPrice:        r.GasPrice,
			IsError:         r.IsError,
			TxReceiptStatus: r.TxReceiptStatus,
		}

		out = append(out, tx)
	}

	var highestBlock string
	if len(out) > 0 {
		highestBlock = out[len(out)-1].BlockNumber
	} else {
		highestBlock = lastBlock
	}

	return &FetchTransactionsResult{
		Transactions: out,
		LastBlock:    highestBlock,
	}, nil
}
