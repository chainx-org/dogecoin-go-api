package main

import "log"
import "github.com/chainx-org/dogecoin-go-api/pkg/dogecoin"

func main() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.SetPrefix("[Dogecoin]")

	secret0 := "flame flock chunk trim modify raise rough client coin busy income smile"
	secret1 := "shrug argue supply evolve alarm caught swamp tissue hollow apology youth ethics"
	secret2 := "awesome beef hill broccoli strike poem rebel unique turn circle cool system"
	priv0, err := dogecoin.GenerateMyPrivkey(secret0)
	if err != nil {
		log.Fatal(err)
	}
	priv1, err := dogecoin.GenerateMyPrivkey(secret1)
	if err != nil {
		log.Fatal(err)
	}
	priv2, err := dogecoin.GenerateMyPrivkey(secret2)
	if err != nil {
		log.Fatal(err)
	}
	pubkey0, err := dogecoin.GenerateMyPubkey(priv0)
	if err != nil {
		log.Fatal(err)
	}
	pubkey1, err := dogecoin.GenerateMyPubkey(priv1)
	if err != nil {
		log.Fatal(err)
	}
	pubkey2, err := dogecoin.GenerateMyPubkey(priv2)
	if err != nil {
		log.Fatal(err)
	}

	redeem_script, err := dogecoin.GenerateRedeemScript([]string{pubkey0, pubkey1, pubkey2}, 2)
	if err != nil {
		log.Fatal(err)
	}
	address, err := dogecoin.GenerateMultisigAddress(redeem_script, "testnet")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("address: %s\n", address)
}
