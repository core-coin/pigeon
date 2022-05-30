package cmd

import (
	"fmt"

	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/core-coin/go-core/accounts/keystore"
	"github.com/core-coin/go-core/common"
	"github.com/core-coin/go-core/common/hexutil"
	"github.com/core-coin/go-goldilocks"

	"github.com/core-coin/pigeon/domain"
	"github.com/core-coin/pigeon/infrastructure/rpcClient/gocore"
	"github.com/core-coin/pigeon/logger/zap"
	txlistuc "github.com/core-coin/pigeon/transaction_list/usecase"
)

func execute() {
	var (
		privateKey *goldilocks.PrivateKey
		err        error
	)

	if verbosityFlag > 7 {
		verbosityFlag = 7
	}
	logger := zap.NewApiLogger(verbosityFlag)
	logger.InitLogger()

	common.DefaultNetworkID = common.NetworkID(networkIDFlag)

	if privateKeyFileFlag != "" && UTCFileFlag != "" {
		logger.Fatal("Cannot use both fags for private key and encrypted UTC file")
		return
	}
	if privateKeyFileFlag != "" {
		privateKey, err = getPrivateKey(privateKeyFileFlag)
		if err != nil {
			logger.Fatalf("Error on getting private key from file: %v", err)
			return
		}
	}
	if UTCFileFlag != "" {
		privateKey, err = getPrivateKeyFromUTC(UTCFileFlag, UTCFilePasswordFlag)
		if err != nil {
			logger.Fatalf("Error on getting private key from UTC file: %v", err)
			return
		}
	}
	rpcClient := gocore.NewRPCClient(gocoreAddressFlag, time.Second*5)
	uc := txlistuc.NewTransactionListUsecase(rpcClient, logger)

	// Get signed transactions from file and stream them
	{
		if signedTxFileFlag != "" {
			txList, err := uc.GetSignedTxsFromFile(signedTxFileFlag)
			if err != nil {
				logger.Fatalf("Error on getting signed transactions from file: %v", err)
			}
			logger.Infof("Successfully got signed transactions from file %v", signedTxFileFlag)
			if !dryrunFlag {
				txIDs, err := uc.StreamSignedTxs(txList)
				if err != nil {
					logger.Errorf("Error on streaming transactions to blockchain: %v", err)
					if len(txIDs) > 0 {
						logger.Error("But some transactions were streamed before error:")
						for i, txID := range txIDs {
							logger.Errorf("%v: %v", i+1, txID)
						}
					}
					return
				}
				logger.Info("Successfully streamed signed transactions into blockchain")

				err = exportTxIDs(uc, txIDs, signedTxResultFileFlag)
				if err != nil {
					logger.Fatalf("Error on exporting transaction hashes: %v", err)
				}
			} else {
				logger.Info("Transactions were not streamed because of dry run!")
			}
			return
		}
	}

	//Get transactions from file, sign them and then stream
	{
		// Get transactions
		txList, err := uc.GetTxsFromFile(txFileFlag, titlesFlag)
		if err != nil {
			logger.Fatalf("Error on getting transactions from file: %v", err)
		}
		logger.Infof("Successfully got transactions from file %v", txFileFlag)
		// Sign transactions
		signedTxs, err := uc.SignTxs(txList, privateKey)
		if err != nil {
			logger.Fatalf("Error on signing transactions from file: %v", err)
		}
		logger.Info("Successfully signed transactions")

		// Save signed transactions into a file if needed
		if exportTxFileFlag != "" {
			err = uc.WriteSignedTxsToFile(signedTxs, exportTxFileFlag)
			if err != nil {
				logger.Fatalf("Error on writing signed transactions to file: %v", err)
			}
			logger.Infof("Successfully saved signed transactions into a file %v", exportTxFileFlag)
			return
		}

		// Stream signed transactions
		if !dryrunFlag {
			txIDs, err := uc.StreamSignedTxs(signedTxs)
			if err != nil {
				logger.Errorf("Error on streaming transactions to blockchain: %v", err)
				if len(txIDs) > 0 {
					logger.Error("But some transactions were streamed before error:")
					for i, txID := range txIDs {
						logger.Errorf("%v: %v", i+1, txID)
					}
				}
				return
			}
			logger.Info("Successfully streamed signed transactions into blockchain")
			err = exportTxIDs(uc, txIDs, signedTxResultFileFlag)
			if err != nil {
				logger.Fatalf("Error on exporting transaction hashes: %v", err)
			}
		} else {
			logger.Info("Transactions were not streamed because of dry run!")
		}
	}
}

func exportTxIDs(uc domain.TransactionListUseCase, txIDs []string, exportFile string) error {
	var err error
	if exportFile != "" {
		err = uc.WriteTxIDsToFile(txIDs, exportFile)
	} else {
		err = uc.WriteTxIDsToConsole(txIDs)
	}
	return err
}

func getPrivateKey(fileName string) (*goldilocks.PrivateKey, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	hexData, err := hexutil.Decode(string(data))
	if err != nil {
		return nil, err
	}
	prKey := goldilocks.BytesToPrivateKey(hexData)
	return &prKey, nil
}

func getPrivateKeyFromUTC(UTCFileName, UTCPasswordFileName string) (*goldilocks.PrivateKey, error) {
	var password string

	jsonBlob, err := os.ReadFile(UTCFileName)
	if err != nil {
		return nil, err
	}

	if UTCPasswordFileName == "" {
		fmt.Print("Enter password for UTC file: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, err
		}
		password = string(bytePassword)
	} else {
		bytePassword, err := os.ReadFile(UTCPasswordFileName)
		if err != nil {
			return nil, err
		}
		password = string(bytePassword)
	}
	key, err := keystore.DecryptKey(jsonBlob, strings.TrimSpace(password))
	if err != nil {
		return nil, err
	}
	return key.PrivateKey, nil
}
