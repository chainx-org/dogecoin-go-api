# Construct Dogecoin Transaction

## Functions

**func BuildTx(baseTx string, signature string, txid string, index uint32, sigType uint32, script string) (string, error)**
**BuildTx**

```
    @Description: Combining signatures into transaction
    @param baseTx: base transaction hex string
    @param signature: signature of sighash
    @param txid: utxo's txid
    @param index: utxo's index
    @param sigType: support [0, 1]. 0 is p2kh, 1 is p2sh
    @param script: When p2kh, script input user pubkey, when p2sh script input redeem script
    @return string: base transaction with one more signature
    @return error: nil|errMsg
```

**func GenerateAddress(pubkey string, network string) (string, error)
GenerateAddress**

```
    @Description: Generate dogecoin p2kh address from pubkey
    @param pubkey: pubkey hex string
    @param network: network string, support  ["mainnet", "testnet"]
    @return string: the dogecoin address string
    @return error: nil|errMsg
```

**func GenerateMultisigAddress(redeemScript string, network string) (string, error)
GenerateMultisigAddress**

```
    @Description: Generate dogecoin p2sh address
    @param redeemScript: redeem script
    @param network: network string, support ["mainnet", "testnet"]
    @return string: the dogecoin address string
    @return error: nil|errMsg
```

**func GenerateMyPrivkey(words string) (string, error)
GenerateMyPrivkey**

```
    @Description: Generate private key from mnemonic
    @param words: mnemonic
    @return string: the private key hex string
    @return error: nil|errMsg
```

**func GenerateMyPubkey(privkey string) (string, error)
GenerateMyPubkey**

```
    @Description: Generate pubkey from privkey
    @param privkey: private key
    @return string: pubkey hex string
    @return error: nil|errMsg
```

**func GenerateRawTransaction(txids []string, indexs []uint32, addresses []string, amounts []uint64) (string, error)
GenerateRawTransaction**

```
    @Description: Add the first input to initialize basic transactions
    @param txids: utxo's txid array
    @param indexs: utxo's index array
    @param addresses: utxo's addresse array
    @param amounts: utxo's amount array
    @return string: the dogecoin raw tx hex string without signature
    @return error: nil|errMsg
```

**func GenerateRedeemScript(pubkeys []string, threshold uint32) (string, error)
GenerateRedeemScript**

```
    @Description: Generate redeem script
    @param pubkeys: hex string concatenated with multiple pubkeys
    @param threshold: threshold number
    @return string: the dogecoin redeem script
    @return error: nil|errMsg
```

**func GenerateSighash(baseTx string, txid string, index uint32, sigType uint32, script string) (string, error)
GenerateSighash**

```
    @Description: Generate sighash/message to sign.Through sigType and script input different, support p2kh and p2sh two types
									of sighash
    @param baseTx: base transaction hex string
    @param txid: utxo's txid
    @param index: utxo's index
    @param sigType: support  [0, 1]. 0 is p2kh, 1 is p2sh
    @param script: When p2kh, script input user pubkey, when p2sh script input redeem script
    @return string: the sighash hex string
    @return error: nil|errMsg
```

**func GenerateSignature(message string, privkey string) (string, error)
GenerateSignature**

```
    @Description: Generate ecdsa signature
    @param message: Awaiting signed sighash/message
    @param privkey: private key
    @return string: the signature hex string
    @return error: nil|errMsg
```

## Example

Here is an [example](https://github.com/chainx-org/dogecoin-go-api/blob/main/demo/dogecoinDemo.go) for reference, which is easy to understand the complete process of using dogecoin package to construct dogecoin transactions.

### Base

1.Convert the mnemonic phrase into a private key

```go
secret0 := "flame flock chunk trim modify raise rough client coin busy income smile"
priv0, err := dogecoin.GenerateMyPrivkey(secret0)
```

2.Get the public key corresponding to the private key

```go
pubkey0, err := dogecoin.GenerateMyPubkey(priv0)
```

3.Get dogecoin address (network support ["mainnet", "testnet"])

```go
addr0, err := dogecoin.GenerateAddress(pubkey0, "testnet")
```

### Get multisig address

1.Get redeem script through all trustees’ public key and threshold

```go
redeem_script, err := dogecoin.GenerateRedeemScript([]string{pubkey0, pubkey1, pubkey2}, 2)
```

2.Generate multisig address by redeem script  (network support ["mainnet", "testnet"])

```go
mutliAddress, err := dogecoin.GenerateMultisigAddress(redeem_script, "testnet")
```

### Construct p2pkh transaction

P2PKH is used for ordinary Dogecoin transfer.

1.Constructing unsigned transactions, you need pass unspent transaction ids, unspent transaction indexs, output addresses and output amounts

```go
base_tx, err := dogecoin.GenerateRawTransaction([]string{"55728d2dc062a9dfe21bae44e87665b270382c8357f14b2a1a4b2b9af92a894a"}, []uint32{0}, []string{addr0, op_return, addr1}, []uint64{100000, 0, 800000})
```

2.Calculate the signature hash for each output (this type of transaction signature is 0)

```go
sighash, err := dogecoin.GenerateSighash(base_tx, txids[i], indexs[i], 0, pubkey1)
```

3.Use the private key to sign for all inputs of transaction

```go
signature, err := dogecoin.GenerateSignature(sighash, priv1)
```

4.Put all the signatures into the transaction

```go
final_tx, err = dogecoin.BuildTx(base_tx, signature, txids[i], indexs[i], 0, pubkey1)
```

5.After all the signatures are put into the base transaction, the transaction structure is completed

### Construct p2sh transaction

P2SH is used for multi -signature transfer.

1.Constructing unsigned transactions, you need pass unspent transaction ids, unspent transaction indexs, output addresses and output amounts

```go
base_tx, err = dogecoin.GenerateRawTransaction([]string{"55728d2dc062a9dfe21bae44e87665b270382c8357f14b2a1a4b2b9af92a894a"}, []uint32{1}, []string{addr1, mutliAddress}, []uint64{1000000, 6000000})
```

2.Calculate the signature hash for each input (Note that you need to use redeem script to calculate the sighash, the signature type is 1)

```go
sighash, err := dogecoin.GenerateSighash(base_tx, txids[i], indexs[i], 1, redeem_script)
```

3.The trustee signed all inputs of the transaction

```go
signature1, err := dogecoin.GenerateSignature(sighash, priv1)
```

4.The trustee put all the signature in the transaction

```go
base_tx, err = dogecoin.BuildTx(base_tx, signature1, txids[i], indexs[i], 1, redeem_script)
```

5.When the number of signatures reaches the threshold, the transaction structure is completed