package cmd

import (
	"errors"
	"log"

	"github.com/dgraph-io/badger"
	"github.com/spf13/cobra"
	"github.com/tjper/golang-blockchain/blockchain"
)

func init() {
	rootCmd.AddCommand(printCmd)
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print the blockchain.",
	Long:  "Print each block in the block chain and its details.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("no arguments expected for print")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var db, err = badger.Open(badger.DefaultOptions(badgerDir))
		if err != nil {
			log.Panic(err)
		}
		defer db.Close()

		var bc = blockchain.New(
			blockchain.WithBadgerDB(db),
		)
		bc.Init()
		bc.Print()
	},
}
