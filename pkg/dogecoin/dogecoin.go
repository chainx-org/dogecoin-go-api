package dogecoin

/*
#cgo android,arm64 LDFLAGS: -L../../lib/android/arm64-v8a -ldogecoin_dll -lm -ldl
#cgo android,arm LDFLAGS: -L../../lib/android/armeabi-v7a -ldogecoin_dll -lm -ldl
#cgo ios LDFLAGS: -L../../lib/ios -ldogecoin_dll -framework Foundation -framework Security
#cgo darwin LDFLAGS: -L../../lib/darwin -ldogecoin_dll
#cgo linux LDFLAGS: -L../../lib/linux -ldogecoin_dll
#include <stdlib.h>
#include "../../lib/DogecoinHeader.h"
*/
import "C"

import (
	"encoding/hex"
	"errors"
	"strings"
	"unsafe"
)

func verifyResult(result *C.char) (string, error) {
	output := C.GoString(result)
	_, err := hex.DecodeString(output)
	if err != nil {
		return "", errors.New(output)
	} else {
		return output, nil
	}
}

// GenerateMyPrivkey
//  @Description: Generate private key from mnemonic
//  @param words: mnemonic
//  @return string: the private key hex string
//  @return error: nil|errMsg
func GenerateMyPrivkey(words string) (string, error) {
	cWords := C.CString(words)
	defer C.free(unsafe.Pointer(cWords))
	salt := C.CString("")
	defer C.free(unsafe.Pointer(salt))
	o := C.generate_my_privkey_dogecoin(cWords, salt)
	return verifyResult(o)
}

// GenerateMyPubkey
//  @Description: Generate pubkey from privkey
//  @param privkey: private key
//  @return string: pubkey hex string
//  @return error: nil|errMsg
func GenerateMyPubkey(privkey string) (string, error) {
	cPrivkey := C.CString(privkey)
	defer C.free(unsafe.Pointer(cPrivkey))
	o := C.generate_my_pubkey_dogecoin(cPrivkey)
	return verifyResult(o)
}

// GenerateAddress
//  @Description: Generate dogecoin p2kh address from pubkey
//  @param pubkey: pubkey hex string
//  @param network: network string, support  ["mainnet", "testnet"]
//  @return string: the dogecoin address string
//  @return error: nil|errMsg
func GenerateAddress(pubkey string, network string) (string, error) {
	cPubkey := C.CString(pubkey)
	defer C.free(unsafe.Pointer(cPubkey))
	cNetwork := C.CString(network)
	defer C.free(unsafe.Pointer(cNetwork))
	o := C.generate_address_dogecoin(cPubkey, cNetwork)
	return C.GoString(o), nil
}

// GenerateRedeemScript
//  @Description: Generate redeem script
//  @param pubkeys: hex string concatenated with multiple pubkeys
//  @param threshold: threshold number
//  @return string: the dogecoin redeem script
//  @return error: nil|errMsg
func GenerateRedeemScript(pubkeys []string, threshold uint32) (string, error) {
	allPubkeys := strings.Join(pubkeys, "")
	cAllPubkeys := C.CString(allPubkeys)
	defer C.free(unsafe.Pointer(cAllPubkeys))
	o := C.generate_redeem_script_dogecoin(cAllPubkeys, C.uint32_t(threshold))
	return verifyResult(o)
}

// GenerateMultisigAddress
//  @Description: Generate dogecoin p2sh address
//  @param redeemScript: redeem script
//  @param network: network string, support ["mainnet", "testnet"]
//  @return string: the dogecoin address string
//  @return error: nil|errMsg
func GenerateMultisigAddress(redeemScript string, network string) (string, error) {
	cRedeemScript := C.CString(redeemScript)
	defer C.free(unsafe.Pointer(cRedeemScript))
	cNetwork := C.CString(network)
	defer C.free(unsafe.Pointer(cNetwork))
	o := C.generate_multisig_address_dogecoin(cRedeemScript, cNetwork)
	return C.GoString(o), nil
}

// GenerateRawTransaction
//  @Description: Add the first input to initialize basic transactions
//  @param txids: utxo's txid array
//  @param indexs: utxo's index array
//  @param addresses: utxo's addresse array
//  @param amounts: utxo's amount array
//  @return string: the dogecoin raw tx hex string without signature
//  @return error: nil|errMsg
func GenerateRawTransaction(txids []string, indexs []uint32, addresses []string, amounts []uint64) (string, error) {
	if len(txids) != len(indexs) || len(addresses) != len(amounts) {
		return "", errors.New("input or output must be equal")
	}

	cTxid := C.CString(txids[0])
	defer C.free(unsafe.Pointer(cTxid))
	baseTx := C.generate_base_tx_dogecoin(cTxid, C.uint32_t(indexs[0]))

	for i := 1; i < len(txids); i++ {
		cTxid = C.CString(txids[i])
		baseTx = C.add_input_dogecoin(baseTx, cTxid, C.uint32_t(indexs[i]))
	}

	var cAddress *C.char
	defer C.free(unsafe.Pointer(cAddress))
	for i := 0; i < len(addresses); i++ {
		cAddress = C.CString(addresses[i])
		baseTx = C.add_output_dogecoin(baseTx, cAddress, C.uint64_t(amounts[i]))
	}

	return verifyResult(baseTx)
}

// GenerateSighash
//  @Description: Generate sighash/message to sign.
//
// Through sigType and script input different, support p2kh and p2sh two types of sighash
//
//  @param baseTx: base transaction hex string
//  @param txid: utxo's txid
//  @param index: utxo's index
//  @param sigType: support  [0, 1]. 0 is p2kh, 1 is p2sh
//  @param script: When p2kh, script input user pubkey, when p2sh script input redeem script
//  @return string: the sighash hex string
//  @return error: nil|errMsg
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

// GenerateSignature
//  @Description: Generate ecdsa signature
//  @param message: Awaiting signed sighash/message
//  @param privkey: private key
//  @return string: the signature hex string
//  @return error: nil|errMsg
func GenerateSignature(message string, privkey string) (string, error) {
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	cPrivkey := C.CString(privkey)
	defer C.free(unsafe.Pointer(cPrivkey))
	o := C.generate_signature_dogecoin(cMessage, cPrivkey)
	return verifyResult(o)
}

// BuildTx
//  @Description: Combining signatures into transaction
//  @param baseTx: base transaction hex string
//  @param signature: signature of sighash
//  @param txid: utxo's txid
//  @param index: utxo's index
//  @param sigType: support [0, 1]. 0 is p2kh, 1 is p2sh
//  @param script: When p2kh, script input user pubkey, when p2sh script input redeem script
//  @return string: base transaction with one more signature
//  @return error: nil|errMsg
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
