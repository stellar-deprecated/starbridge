import StellarSdk from 'stellar-sdk'

const server = new StellarSdk.Server('https://horizon-testnet.stellar.org')

const createPaymentTransaction = async (
  publicKey: string,
  publicKeyDestination: string,
  amount: string
): Promise<void> => {
  const account = await server.loadAccount(publicKey)

  /*
        Right now, we have one function that fetches the base fee.
        In the future, we'll have functions that are smarter about suggesting fees,
        e.g.: `fetchCheapFee`, `fetchAverageFee`, `fetchPriorityFee`, etc.
    */
  const fee = await server.fetchBaseFee()

  const transaction = new StellarSdk.TransactionBuilder(account, {
    fee,
    networkPassphrase: StellarSdk.Networks.TESTNET,
  })
    .addOperation(
      // this operation funds the new account with XLM
      StellarSdk.Operation.payment({
        destination: publicKeyDestination,
        asset: StellarSdk.Asset.native(),
        amount: amount,
      })
    )
    .setTimeout(30)
    .build()

  // sign the transaction
  //   transaction.sign(StellarSdk.Keypair.fromSecret(secretString))

  try {
    const transactionResult = await server.submitTransaction(transaction)
    console.log(transactionResult)
  } catch (err) {
    console.error(err)
  }
}

export { createPaymentTransaction }
