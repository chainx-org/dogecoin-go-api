#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

/**
 * Generate private key from mnemonic
 *
 * - [`phrase`]: root phrase
 * - [`pd_passphrase`]: pass phrase
 *
 * Return the private key hex string.
 */
char *generate_my_privkey_dogecoin(const char *phrase, const char *pd_passphrase);

/**
 * Generate pubkey from privkey
 *
 * - [`privkey`]: private key
 *
 * Returns: pubkey string
 */
char *generate_my_pubkey_dogecoin(const char *privkey);

/**
 * Generate dogecoin p2kh address from pubkey
 *
 * - [`pubkey`]: pubkey hex string
 * - [`network`]: network string, support [`mainnet` ,`testnet`]
 *
 * Return the dogecoin address string.
 */
char *generate_address_dogecoin(const char *pubkey, const char *network);

/**
 * Generate redeem script
 *
 * - [`pubkeys`]: Hex string concatenated with multiple pubkeys
 * - [`threshold`]: threshold number
 *
 * Return the dogecoin redeem script.
 */
char *generate_redeem_script_dogecoin(const char *pubkeys, uint32_t threshold);

/**
 * Generate dogecoin p2sh address
 *
 * - [`redeem_script`]: redeem script
 * - [`network`]: network string, support [`mainnet` ,`testnet`]
 *
 * Return the dogecoin address string.
 */
char *generate_multisig_address_dogecoin(const char *redeem_script, const char *network);

/**
 * Add the first input to initialize basic transactions
 *
 * - [`txid`]: utxo's txid
 * - [`index`]: utxo's index
 *
 * Return the base tx hex string with one input.
 */
char *generate_base_tx_dogecoin(const char *txid, uint32_t index);

/**
 * Increase the input with txid and index
 *
 * - [`base_tx`]: base transaction hex string
 * - [`txid`]: utxo's txid
 * - [`index`]: utxo's index
 *
 * Return the base tx hex string with more input.
 */
char *add_input_dogecoin(const char *base_tx, const char *txid, uint32_t index);

/**
 * Increase the output with address and amount
 *
 * - [`base_tx`]: base transaction hex string
 * - [`address`]: output address
 * - [`amount`]: output amount
 *
 * Return the base tx hex string with more output.
 */
char *add_output_dogecoin(const char *base_tx, const char *address, uint64_t amount);

/**
 * Generate sighash/message to sign.
 *
 * NOTE: Through sig_type and script input different, support p2kh and p2sh two types of sighash.
 *
 * - [`base_tx`]: base transaction hex string
 * - [`txid`]: utxo's txid
 * - [`index`]: utxo's index
 * - [`sig_type`]: support [`0`, `1``], 0 is p2kh, 1 is p2sh
 * - [`script`]: When p2kh, script input user pubkey, when p2sh script input redeem script
 *
 * Return the sighash string.
 */
char *generate_sighash_dogecoin(const char *base_tx,
                                const char *txid,
                                uint32_t index,
                                uint32_t sig_type,
                                const char *script);

/**
 * Generate ecdsa signature.
 *
 * [`message`]: Awaiting signed sighash/message
 * [`privkey`]: private key
 *
 * Return the signature hex string.
 */
char *generate_signature_dogecoin(const char *message, const char *privkey);

/**
 * Combining signatures into transaction.
 *
 * - [`base_tx`]: base transaction hex string.
 * - [`signature`]: signature of sighash
 * - [`txid`]: utxo's txid
 * - [`index`]: utxo's index
 * - [`sig_type`]: support [`0`, `1``], 0 is p2kh, 1 is p2sh
 * - [`script`]: When p2kh, script input user pubkey, when p2sh script input redeem script
 *
 * Return transaction with one more signature.
 */
char *build_tx_dogecoin(const char *base_tx,
                        const char *signature,
                        const char *txid,
                        uint32_t index,
                        uint32_t sig_type,
                        const char *script);
