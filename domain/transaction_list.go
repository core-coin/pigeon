package domain

import (
	"github.com/core-coin/go-goldilocks"
)

type TransactionList []*Transaction

type Transaction struct {
	To          string  `json:"to" csv:"to"`
	From        string  `json:"from" csv:"from"`
	Amount      float64 `json:"amount" csv:"amount"`
	EnergyLimit string  `json:"energy_limit" csv:"energy_limit"`
	EnergyPrice string  `json:"energy_price" csv:"energy_price"`
	Nonce       string  `json:"nonce" csv:"nonce"`
}

type TransactionListUseCase interface {
	//StreamSignedTxs is receiving a file with signed transactions and stream them into a blockchain
	// Returns a slice of IDs of sent transactions
	StreamSignedTxs(signedTxs []string) ([]string, error)
	//WriteTxIDsToFile is receiving a slice of transaction IDs and write them to a file
	WriteTxIDsToFile(txIDs []string, fileName string) error
	//WriteTxIDsToConsole is receiving a slice of transaction IDs and write them to a console
	WriteTxIDsToConsole(txIDs []string) error
	//GetSignedTxsFromFile is reading signed transactions from a file
	GetSignedTxsFromFile(fileName string) ([]string, error)
	//GetTxsFromFile is reading transaction from a file and skip first row in CSV if missTitles is true
	GetTxsFromFile(fileName string, missTitles bool) (TransactionList, error)
	//SignTxs signs transactions with provided private key
	SignTxs(txs TransactionList, key *goldilocks.PrivateKey) ([]string, error)
	//WriteSignedTxsToFile is writing signed transactions into a file in JSON format
	WriteSignedTxsToFile(signedTxs []string, fileName string) error
}
