package main

import (
	"log"

	"github.com/chainx-org/dogecoin-go-api/pkg/dogecoin"
)

func main() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.SetPrefix("[Dogecoin]")
	// 1. phrase -> private -> pubkey
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

	// 2. Generate p2kh address
	addr0, err := dogecoin.GenerateAddress(pubkey0, "testnet")
	if err != nil {
		log.Fatal(err)
	}
	addr1, err := dogecoin.GenerateAddress(pubkey1, "testnet")
	if err != nil {
		log.Fatal(err)
	}

	// 3. Generate p2sh address
	// step0: generate redeem script using multiple user public keys and threshold
	redeem_script, err := dogecoin.GenerateRedeemScript([]string{pubkey0, pubkey1, pubkey2}, 2)
	if err != nil {
		log.Fatal(err)
	}
	// setp1: Use reddem script to generate addresses
	mutliAddress, err := dogecoin.GenerateMultisigAddress(redeem_script, "testnet")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("multiAddress:", mutliAddress)

	// 4. spent from p2kh address, e.g: addr1
	// step0: enter txid, index, address, amount array to generate raw tx
	op_return := "35516a706f3772516e7751657479736167477a6334526a376f737758534c6d4d7141754332416255364c464646476a38"
	txids := []string{"55728d2dc062a9dfe21bae44e87665b270382c8357f14b2a1a4b2b9af92a894a"}
	indexs := []uint32{0}
	base_tx, err := dogecoin.GenerateRawTransaction(txids, indexs, []string{addr0, op_return, addr1}, []uint64{100000, 0, 800000})
	if err != nil {
		log.Fatal(err)
	}

	var final_tx string
	// note: Each txid must be signed separately
	for i := 0; i < len(txids); i++ {
		// step1: generate p2kh sighash
		sighash, err := dogecoin.GenerateSighash(base_tx, txids[i], indexs[i], 0, pubkey1)
		if err != nil {
			log.Fatal(err)
		}
		// step2: generate signature
		signature, err := dogecoin.GenerateSignature(sighash, priv1)
		if err != nil {
			log.Fatal(err)
		}
		// step3: combine base tx and signature
		final_tx, err = dogecoin.BuildTx(base_tx, signature, txids[i], indexs[i], 0, pubkey1)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("final p2kh tx: %s\n", final_tx)

	// 5. spent from p2sh address, e.g: multi_addr
	// step0: enter txid, index, address, amount array to generate raw tx
	txids = []string{"55728d2dc062a9dfe21bae44e87665b270382c8357f14b2a1a4b2b9af92a894a"}
	indexs = []uint32{1}
	base_tx, err = dogecoin.GenerateRawTransaction(txids, indexs, []string{addr1, mutliAddress}, []uint64{1000000, 6000000})
	if err != nil {
		log.Fatal(err)
	}
	// note: Each txid must be signed separately, and each tixd has at least a threshold of user signatures!
	for i := 0; i < len(txids); i++ {
		// step1: generate p2sh sighash
		sighash, err := dogecoin.GenerateSighash(base_tx, txids[i], indexs[i], 1, redeem_script)
		if err != nil {
			log.Fatal(err)
		}
		// step2: user1 generate signature
		signature1, err := dogecoin.GenerateSignature(sighash, priv1)
		if err != nil {
			log.Fatal(err)
		}
		// step3: combine base tx and signature
		base_tx, err = dogecoin.BuildTx(base_tx, signature1, txids[i], indexs[i], 1, redeem_script)
		if err != nil {
			log.Fatal(err)
		}
		// step4: user0 generate signature
		signature0, err := dogecoin.GenerateSignature(sighash, priv0)
		if err != nil {
			log.Fatal(err)
		}
		// step5: combine base tx and signature
		final_tx, err = dogecoin.BuildTx(base_tx, signature0, txids[i], indexs[i], 1, redeem_script)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("final p2sh tx: %s\n", final_tx)
}
