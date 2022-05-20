package usecase

import (
	"encoding/json"
	"errors"
	"github.com/core-coin/tx-signer/pkg"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"

	"github.com/core-coin/go-core/common"
	"github.com/core-coin/go-core/common/hexutil"
	"github.com/core-coin/go-core/core/types"
	"github.com/core-coin/go-core/rlp"
	"github.com/core-coin/go-goldilocks"
	"github.com/gocarina/gocsv"

	"github.com/core-coin/tx-signer/domain"
	"github.com/core-coin/tx-signer/infrastructure/rpcClient"
	"github.com/core-coin/tx-signer/logger"
)

type transactionListUsecase struct {
	logger logger.Logger
	rpc    rpcClient.Client
}

//NewTransactionListUsecase create new transaction list usecase
func NewTransactionListUsecase(rpc rpcClient.Client, log logger.Logger) domain.TransactionListUseCase {
	return &transactionListUsecase{
		rpc:    rpc,
		logger: log,
	}
}

//StreamSignedTxs is sending raw transactions to blockchain
func (t *transactionListUsecase) StreamSignedTxs(signedTxs []string) ([]string, error) {
	var txIDs []string
	for _, tx := range signedTxs {
		hash, err := t.rpc.SendRawTransaction(tx)
		if err != nil {
			return txIDs, err
		}
		t.logger.Debugf("Streamed transaction with hash %v", hash)
		txIDs = append(txIDs, hash)
	}
	return txIDs, nil
}

//WriteTxIDsToFile is writing transaction hashes to file
func (t *transactionListUsecase) WriteTxIDsToFile(txIDs []string, fileName string) error {
	if len(txIDs) == 0 {
		t.logger.Debug("Trying to write 0 txs IDs to file")
		return nil
	}

	data, err := json.Marshal(txIDs)
	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}
	t.logger.Infof("Transactions hashes were saved to file %v", fileName)
	return nil
}

//WriteTxIDsToConsole is writing transaction hashes to console
func (t *transactionListUsecase) WriteTxIDsToConsole(txIDs []string) error {
	if len(txIDs) == 0 {
		t.logger.Debug("Trying to write 0 txs IDs to console")
		return nil
	}
	if len(txIDs) == 1 {
		t.logger.Infof("Transaction %v was streamed successfully", txIDs[0])
		return nil
	}
	t.logger.Info("Transactions were streamed successfully:")
	for i, tx := range txIDs {
		t.logger.Infof("%v: %v", i, tx)
	}
	return nil
}

//GetSignedTxsFromFile is getting raw transaction from file
func (t *transactionListUsecase) GetSignedTxsFromFile(fileName string) ([]string, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return []string{}, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []string{}, err
	}

	var result []string
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

//GetTxsFromFile is getting transactions from file
func (t *transactionListUsecase) GetTxsFromFile(fileName string, missTitles bool) (domain.TransactionList, error) {
	txsFromFile, err := t.getTxsFromFile(fileName, missTitles)
	if err != nil {
		panic(err)
	}

	nonces := map[string]string{}

	for _, tx := range txsFromFile {
		// set default to empty values
		if tx.Nonce == "" {
			if nonce, ok := nonces[tx.From]; !ok {
				nonce, err := t.rpc.GetAccountNonce(tx.From, "pending")
				if err != nil {
					if err != nil {
						panic(err)
					}
				}
				tx.Nonce = strconv.FormatUint(nonce, 10)
				nonces[tx.From] = tx.Nonce
			} else {
				nonceInt, err := strconv.Atoi(nonce)
				if err != nil {
					if err != nil {
						panic(err)
					}
				}
				tx.Nonce = strconv.FormatUint(uint64(nonceInt+1), 10)
				nonces[tx.From] = tx.Nonce
			}
		}
		if tx.EnergyPrice == "" {
			energyPrice, err := t.rpc.EstimateEnergyPrice()
			if err != nil {
				panic(err)
			}
			tx.EnergyPrice = strconv.Itoa(int(energyPrice))
		}
		if tx.EnergyLimit == "" {
			tx.EnergyLimit = "21000"
		}

	}
	return txsFromFile, nil
}

//SignTxs signs transactions
func (t *transactionListUsecase) SignTxs(txs domain.TransactionList, key *goldilocks.PrivateKey) ([]string, error) {
	var signed []string
	signer := types.MakeSigner(big.NewInt(int64(common.DefaultNetworkID)))
	for _, internalTx := range txs {
		tx, err := t.TxToGocoreType(internalTx)
		if err != nil {
			return signed, err
		}

		signedTx, err := types.SignTx(tx, signer, key)
		if err != nil {
			return signed, err
		}

		signedTxBytes, err := rlp.EncodeToBytes(signedTx)
		if err != nil {
			return signed, err
		}
		signed = append(signed, hexutil.Encode(signedTxBytes))
	}
	return signed, nil
}

//WriteSignedTxsToFile is writing transactions to file
func (t *transactionListUsecase) WriteSignedTxsToFile(signedTxs []string, fileName string) error {
	if len(signedTxs) == 0 {
		t.logger.Debug("Trying to write 0 signed txs to file")
		return nil
	}

	data, err := json.Marshal(signedTxs)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

//getTxsFromFile is loading transactions from file and choose method depending on file extension
func (t *transactionListUsecase) getTxsFromFile(fileName string, missTitles bool) ([]*domain.Transaction, error) {
	switch filepath.Ext(fileName) {
	case ".json":
		return t.getTxsFromJSON(fileName)
	case ".csv":
		return t.getTxsFromCSV(fileName, missTitles)
	}
	return nil, errors.New("unsupported file extension")
}

//getTxsFromJSON is loading transactions from json file
func (t *transactionListUsecase) getTxsFromJSON(fileName string) ([]*domain.Transaction, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var txs []*domain.Transaction
	err = json.Unmarshal(byteValue, &txs)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

//getTxsFromCSV is loading transaction from csv file
func (t *transactionListUsecase) getTxsFromCSV(fileName string, missTitles bool) ([]*domain.Transaction, error) {
	in, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	txs := []*domain.Transaction{}

	read := gocsv.DefaultCSVReader(in)
	if missTitles {
		if err := gocsv.UnmarshalCSVWithoutHeaders(read, &txs); err != nil {
			return nil, err
		}
	} else {
		if err := gocsv.UnmarshalCSV(read, &txs); err != nil {
			return nil, err
		}
	}
	return txs, nil
}

//TxToGocoreType converts *domain.Transaction type to gocore *types.Transaction type
func (t *transactionListUsecase) TxToGocoreType(tx *domain.Transaction) (*types.Transaction, error) {
	to, err := common.HexToAddress(tx.To)
	if err != nil {
		return nil, err
	}

	nonce, err := strconv.Atoi(tx.Nonce)
	if err != nil {
		return nil, err
	}

	limit, err := strconv.Atoi(tx.EnergyLimit)
	if err != nil {
		return nil, err
	}

	price, ok := new(big.Int).SetString(tx.EnergyPrice, 10)
	if !ok {
		return nil, errors.New("energy price in transaction has bad number ")
	}

	bigval := new(big.Float)
	bigval.SetFloat64(tx.Amount)

	coin := new(big.Float)
	coin.SetInt(pkg.Core)

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result)

	gocoreTx := types.NewTransaction(uint64(nonce), to, result, uint64(limit), price, []byte{})
	t.logger.Debugf("Converted transaction from: %+v to %+v", tx, gocoreTx)
	return gocoreTx, nil
}
