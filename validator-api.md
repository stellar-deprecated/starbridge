# User-facing API

## Signature request

Parameters:
* A hash identifying a deposit to the bridge contract on Ethereum

## Refund request

# Stellar-facing API

## Callbacks

* New ledger notification
  * Must include:
    * ledger close time
    * all deposits made to the bridge account in this ledger
    * all withdrawals from the bridge account made in this ledger
  * Precondition: ledger is final

## Calls

* Get the sequence number of an account
* Ingest history (used on startup)
  Provides the history of all withdrawals from the bridge account

# Ethereum-facing API

* New-block notification
  * Must include:
    * block timestamp
    * all deposits made to the bridge contract in this block
    * all withdrawals from the bridge contract in this block
  * Preconditions: block is final (e.g. 6 deep in the chain)

