# Starbridge Protocol

Parameters:
* The *withdrawal window* determines how much time the user has to withdraw the wrapped asset on Stellar.

Bridge validator state:
* `last_ledger`, the last Stellar ledger known to the bridge validator.

## Transferring an Ethereum-native asset to Stellar

### Withdrawing the funds on Stellar

Steps:
1. The user initiates a transfer by sending some count `N` of `T` tokens from their account on Ethereum (the *sending account*) to the bridge account on Ethereum, also specifying the *destination account* on Stellar. A unique *deposit identifier* identifies this operation. Let `t` be the timestamp of the Ethereum block in which this operation is executed.
2. The user creates a *withdraw transaction* and sends a *signature request* to each bridge validator, including the deposit identifier.
3. Every bridge validator does the following:
  1. Check that no withdraw transaction for the same deposit identifier has been executed on Stellar as of ledger `last_ledger`, and return an error if not.
  2. Else, sign the withdraw transaction provided by the user, provided that:
    1. its transferring a count `N` of `wT` tokens, where `wT` is the wrapped counterpart of `T` on Stellar, from the bridge account to the destination account,
    2. it has a time bound of `t` plus the withdraw window,
    3. it has a sequence number of 1 plus the sequence number of the receiving account as of `last_ledger`,
    4. the source account is the destination account,
    5. the memo contains the deposit identifier.
4. Once the user has collected enough validator signatures, it submits the withdraw transaction on Stellar to receive the funds.
5. The withdraw transaction might fail, for example if the sequence number of the receiving account does not correspond to the sequence number of the withdraw transaction.
6. The user can submit signature requests as many times as they want. Bridge validators process each signature request as in point 3 above.

### Cancelling the transfer of an Ethereum-native asset to Stellar

Steps:
1. The user sends a refund request to every bridge validator, providing the deposit identifier. This request must be signed by the key of the sending account.
2. Every bridge validator does the following:
  1. Check that the close time of `last_ledger` is strictly greater than `t` plus the withdrawal window, and return an error if not.
  2. Check that no withdraw transaction for the same deposit identifier has been executed on Stellar as of ledger `last_ledger`, and return an error if not.
  3. Sign a *refund approval* and return it to the user.
3. Once the user has collected enough signed refund approvals, they call the bridge contract on Ethereum to receive their refund.
