package cmd

import (
	"errors"
	"log"

	"github.com/dgraph-io/badger"
	"github.com/spf13/cobra"
	"github.com/tjper/golang-blockchain/blockchain"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

const badgerDir = "/tmp/blocks"

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a block to the blockchain",
	Long:  "By adding block to this blockchain a block is added utilizing Proof of Work consensus algorithm.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires an add argument")
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
		bc.AddBlock(args[0])
	},
}
