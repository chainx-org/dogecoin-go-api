package main

/*
#cgo LDFLAGS: -L./lib -ldogecoin_dll
#include <stdlib.h>
#include "./lib/DogecoinHeader.h"
*/
import "C"

import (
	"encoding/hex"
	"log"
	"strings"
	"unsafe"
)

func verifyResult(result *C.char) (string, error) {
	output := C.GoString(result)
	log.Println(output)
	_, err := hex.DecodeString(output)
	return output, err
}

func GenerateMyPrivkey(words string) (string, error) {
	cWords := C.CString(words)
	defer C.free(unsafe.Pointer(cWords))
	salt := C.CString("")
	defer C.free(unsafe.Pointer(salt))
	o := C.generate_my_privkey_dogecoin(cWords, salt)
	return verifyResult(o)
}

func GenerateMyPubkey(privkey string) (string, error) {
	cPrivkey := C.CString(privkey)
	defer C.free(unsafe.Pointer(cPrivkey))
	o := C.generate_my_pubkey_dogecoin(cPrivkey)
	return verifyResult(o)
}

func GenerateAddress(pubkey string, network string) (string, error) {
	cPubkey := C.CString(pubkey)
	defer C.free(unsafe.Pointer(cPubkey))
	cNetwork := C.CString(network)
	defer C.free(unsafe.Pointer(cNetwork))
	o := C.generate_address_dogecoin(cPubkey, cNetwork)
	return C.GoString(o), nil
}

func GenerateRedeemScript(pubkeys []string, threshold uint32) (string, error) {
	allPubkeys := strings.Join(pubkeys, "")
	cAllPubkeys := C.CString(allPubkeys)
	defer C.free(unsafe.Pointer(cAllPubkeys))
	o := C.generate_redeem_script_dogecoin(cAllPubkeys, C.uint32_t(threshold))
	return verifyResult(o)
}

func GenerateMultisigAddress(redeemScript string, network string) (string, error) {
	cRedeemScript := C.CString(redeemScript)
	defer C.free(unsafe.Pointer(cRedeemScript))
	cNetwork := C.CString(network)
	defer C.free(unsafe.Pointer(cNetwork))
	o := C.generate_multisig_address_dogecoin(cRedeemScript, cNetwork)
	return C.GoString(o), nil
}

func GenerateRawTransaction(txids []string, indexs []uint32, addresses []string, amounts []uint64) (string, error) {
	if len(txids) != len(indexs) || len(addresses) != len(amounts) {
		log.Fatalf("txids, indexs must be same length ,addresses, amounts must be same length\n")
	}

	cTxid0 := C.CString(txids[0])
	defer C.free(unsafe.Pointer(cTxid0))
	base_tx := C.generate_base_tx_dogecoin(cTxid0, C.uint32_t(indexs[0]))

	for i := 1; i < len(txids); i++ {
		cTxid := C.CString(txids[i])
		defer C.free(unsafe.Pointer(cTxid))
		base_tx = C.add_input_dogecoin(base_tx, cTxid, C.uint32_t(indexs[i]))
	}

	for i := 0; i < len(addresses); i++ {
		cAddress := C.CString(addresses[i])
		defer C.free(unsafe.Pointer(cAddress))
		base_tx = C.add_output_dogecoin(base_tx, cAddress, C.uint64_t(amounts[i]))
	}

	return verifyResult(base_tx)
}

func GenerateSighash(baseTx string, txid string, index uint32, sigType uint32, script string) (string, error) {
	cBaseTx := C.CString(baseTx)
	defer C.free(unsafe.Pointer(cBaseTx))
	cTxid := C.CString(txid)
	defer C.free(unsafe.Pointer(cTxid))
	cScript := C.CString(script)
	defer C.free(unsafe.Pointer(cScript))
	o := C.generate_sighash_dogecoin(cBaseTx, cTxid, C.uint32_t(index), C.uint32_t(sigType), cScript)
	return verifyResult(o)
}

func GenerateSignature(message string, privkey string) (string, error) {
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	cPrivkey := C.CString(privkey)
	defer C.free(unsafe.Pointer(cPrivkey))
	o := C.generate_signature_dogecoin(cMessage, cPrivkey)
	return verifyResult(o)
}

func BuildTx(baseTx string, signature string, txid string, index uint32, sigType uint32, script string) (string, error) {
	cBaseTx := C.CString(baseTx)
	defer C.free(unsafe.Pointer(cBaseTx))
	cSignature := C.CString(signature)
	defer C.free(unsafe.Pointer(cSignature))
	cTxid := C.CString(txid)
	defer C.free(unsafe.Pointer(cTxid))
	cScript := C.CString(script)
	defer C.free(unsafe.Pointer(cScript))
	o := C.build_tx_dogecoin(cBaseTx, cSignature, cTxid, C.uint32_t(index), C.uint32_t(sigType), cScript)
	return verifyResult(o)
}

func main() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.SetPrefix("[Dogecoin]")

	secret0 := "flame flock chunk trim modify raise rough client coin busy income smile"
	secret1 := "shrug argue supply evolve alarm caught swamp tissue hollow apology youth ethics"
	secret2 := "awesome beef hill broccoli strike poem rebel unique turn circle cool system"
	priv0, err := GenerateMyPrivkey(secret0)
	if err != nil {
		log.Fatal(err)
	}
	priv1, err := GenerateMyPrivkey(secret1)
	if err != nil {
		log.Fatal(err)
	}
	priv2, err := GenerateMyPrivkey(secret2)
	if err != nil {
		log.Fatal(err)
	}
	pubkey0, err := GenerateMyPubkey(priv0)
	if err != nil {
		log.Fatal(err)
	}
	pubkey1, err := GenerateMyPubkey(priv1)
	if err != nil {
		log.Fatal(err)
	}
	pubkey2, err := GenerateMyPubkey(priv2)
	if err != nil {
		log.Fatal(err)
	}

	redeem_script, err := GenerateRedeemScript([]string{pubkey0, pubkey1, pubkey2}, 2)
	if err != nil {
		log.Fatal(err)
	}
	address, err := GenerateMultisigAddress(redeem_script, "testnet")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("address: %s\n", address)
}
