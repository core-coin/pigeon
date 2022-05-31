package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Flags
var (
	txFileFlag       string
	exportTxFileFlag string

	signedTxFileFlag       string
	signedTxResultFileFlag string

	titlesFlag    bool
	dryrunFlag    bool
	verbosityFlag int

	networkIDFlag       int
	privateKeyFileFlag  string
	UTCFileFlag         string
	UTCFilePasswordFlag string

	gocoreAddressFlag string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "pigeon",
	Example: examples,
	Short:   "Sign & transmit transactions",
	Long:    `This application is used to sign transactions and stream them in Core Blockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		execute()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	RootCmd.PersistentFlags().BoolVarP(&dryrunFlag, "dry-run", "d", false, "Test the schema (do not stream, do not sign)")
	RootCmd.PersistentFlags().BoolVarP(&titlesFlag, "titles", "t", false, "Skip 1 line (for CSV)")
	RootCmd.PersistentFlags().IntVarP(&verbosityFlag, "verbosity ", "v", 2, "Verbosity (from 1 to 7)")

	RootCmd.PersistentFlags().IntVarP(&networkIDFlag, "network", "n", 1, "Network to stream on")
	RootCmd.PersistentFlags().StringVarP(&gocoreAddressFlag, "gocore", "g", "http://127.0.0.1:8545", "Gocore RPC API endpoint")
	RootCmd.PersistentFlags().StringVarP(&privateKeyFileFlag, "private-key-file", "k", "", "File with private key to sign transactions")
	RootCmd.PersistentFlags().StringVarP(&UTCFileFlag, "utc-file", "u", "", "UTC file with encoded private key")
	RootCmd.PersistentFlags().StringVarP(&UTCFilePasswordFlag, "password-file", "p", "", "File with password to for file")

	RootCmd.PersistentFlags().StringVarP(&txFileFlag, "file", "f", "", "Input file with transactions")
	RootCmd.PersistentFlags().StringVarP(&exportTxFileFlag, "output", "o", "", "Output file with signed transactions")

	RootCmd.PersistentFlags().StringVarP(&signedTxFileFlag, "stream-file", "s", "", "File for streaming transactions into blockchain")
	RootCmd.PersistentFlags().StringVarP(&signedTxResultFileFlag, "tx-ids-file", "i", "", "File where to store streamed tx IDs")

}

const examples = `
To sign transactions offline: pigeon -f {path to file with transactions} -u {path to UTC file} -o {path to file where to save signed transactions}
To sign and stream transactions: pigeon -f {path to file with transactions} -u {path to UTC file} -p {path to file with password}
To sign and stream transactions(+ save streamed transaction IDs to file): pigeon -f {path to file with transactions} -u {path to UTC file} -i {path to file where to save transactions hashes}
To stream signed transactions: pigeon -s {path to file with signed transactions}
To stream signed transactions(+ save streamed transaction IDs to file): pigeon -s {path to file with signed transactions} -i {path to file where to save transactions hashes}
`
