package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {

	initMsg := `Hello from stellar bot choose from options below: \n
  1. Generate keypair
  2. Create account
  3. Send payment`

	fmt.Println(initMsg)

	var selectedOption string
	fmt.Print(`Select from options above: `)
	fmt.Scanln(&selectedOption)
	if selectedOption == "1" {
		pair, err := keypair.Random()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Your seed: " + pair.Seed())
		fmt.Println("Your address " + pair.Address())
	} else if selectedOption == "2" {
		var address string
		fmt.Print(`Provide your address: `)
		fmt.Scanln(&address)

		resp, err := http.Get("https://friendbot.stellar.org/?addr=" + address)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))
	} else if selectedOption == "3" {
		var destination string
		fmt.Print(`Provide your destination's address: `)
		fmt.Scanln(&destination)

		var amount string
		fmt.Print(`Provide your amount (in XLM): `)
		fmt.Scanln(&amount)

		var secret string
		fmt.Print(`Provide your secret: `)
		fmt.Scanln(&secret)

		client := horizonclient.DefaultTestNetClient

		// Make sure destination account exists
		destAccountRequest := horizonclient.AccountRequest{AccountID: destination}
		destinationAccount, err := client.AccountDetail(destAccountRequest)
		if err != nil {
			panic(err)
		}

		fmt.Println("Destination Account", destinationAccount)

		// Load the source account
		sourceKP := keypair.MustParseFull(secret)
		sourceAccountRequest := horizonclient.AccountRequest{AccountID: sourceKP.Address()}
		sourceAccount, err := client.AccountDetail(sourceAccountRequest)
		if err != nil {
			panic(err)
		}

		// Build transaction
		tx, err := txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        &sourceAccount,
				IncrementSequenceNum: true,
				BaseFee:              txnbuild.MinBaseFee,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Operations: []txnbuild.Operation{
					&txnbuild.Payment{
						Destination: destination,
						Amount:      amount,
						Asset:       txnbuild.NativeAsset{},
					},
				},
			},
		)

		if err != nil {
			panic(err)
		}

		tx, err = tx.Sign(network.TestNetworkPassphrase, sourceKP)
		if err != nil {
			panic(err)
		}

		resp, err := horizonclient.DefaultTestNetClient.SubmitTransaction(tx)
		if err != nil {
			panic(err)
		}

		fmt.Println("Successful Transaction:")
		fmt.Println("Ledger:", resp.Ledger)
		fmt.Println("Hash:", resp.Hash)
	}

}
