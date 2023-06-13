package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stellar/go/keypair"
)

func main() {

	initMsg := `Hello from stellar bot choose from options below: \n
  1. Generate keypair
  2. Create account`

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
	}

}
