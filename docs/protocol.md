# Starbridge Protocol

Parameters:
* The *withdrawal window* determines how much time the user has to withdraw a wrapped asset on Stellar after making a deposit on Ethereum.

Bridge validator state:
* `last_ledger`, the last Stellar ledger known to the bridge validator.

## Transferring an Ethereum-native asset to Stellar

### Withdrawing the funds on Stellar

Steps:
1. The user initiates a transfer by sending some count `N` of `T` tokens from their account on Ethereum (the *sending account*) to the bridge account on Ethereum, also specifying the *destination account* on Stellar. A unique *deposit identifier* identifies this operation. Let `t` be the timestamp of the Ethereum block in which this operation is executed.
2. The user creates a *withdraw transaction* and sends a *signature request* , including the deposit identifier, to each bridge validator. This request must be signed by the key of the withdrawing account on Stellar.
3. Upon receiving the request, each bridge validator does the following:
  1. Check that no withdraw transaction for the same deposit identifier has been executed on Stellar as of ledger `last_ledger`, and return an error if not.
  2. Else, sign the withdraw transaction provided by the user, provided that:
    1. it's transferring the right count `N` of `wT` tokens, where `wT` is the wrapped counterpart of `T` on Stellar, from the bridge account to the destination account,
    2. it has a time bound of `t` plus the withdraw window,
    3. it has a sequence number of 1 plus the sequence number of the receiving account as of `last_ledger`,
    4. the source account of the transaction is the destination account,
    5. the memo contains the deposit identifier.
4. Once the user has collected enough validator signatures, it submits the withdraw transaction on Stellar to receive the funds.
5. The withdraw transaction might fail, for example if the sequence number of the receiving account does not correspond to the sequence number of the withdraw transaction.
6. The user can submit new signature requests as many times as they want. Bridge validators process each signature request as in point 3 above.

### Cancelling the transfer of an Ethereum-native asset to Stellar

Steps:
1. The user sends a refund request to every bridge validator, providing the deposit identifier. This request must be signed by the key of the sending account.
2. Every bridge validator does the following:
  1. Check that the close time of `last_ledger` is strictly greater than `t` plus the withdrawal window, and return an error if not.
  2. Check that no withdraw transaction for the same deposit identifier has been executed on Stellar as of ledger `last_ledger`, and return an error if not.
  3. Sign a *refund approval* and return it to the user.
3. Once the user has collected enough signed refund approvals, they call the bridge contract on Ethereum to receive their refund.

### Design Rationale

The advantages of this protocol are:
1. Only two transactions needed (deposit on Ethereum and withdraw on Stellar).
2. Bridge validators can always safely restart from a blank state.
3. Safety does not rely on any timing assumptions.

The disadvantages are that it is not very flexible for the user: the withdraw period is fixed and no cancellation is possible before it ends, and the user needs to keep their sequence number stable throughout the process or risk having to create and get signatures for a new withdraw transaction.

## Transferring an Stellar-native asset to Ethereum

TODO

## Formal models

[starbrige-timelock.ivy](./starbridge-timelock.ivy) is a high-level model of the Ethereum to Stellar transfer flow. The model takes into account that validators may lose their state and restart. It includes a proof that no double-spend can ever happen (the proof consists of an inductive invariant).
