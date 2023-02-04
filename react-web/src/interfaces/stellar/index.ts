import { Dispatch } from 'react'
import { toast } from 'react-toastify'

import { arrayify, zeroPad } from '@ethersproject/bytes'
import { signTransaction, getPublicKey } from '@stellar/freighter-api'
import StellarSdk from 'stellar-sdk'

import { validatorUrls, WithdrawResult } from 'interfaces/http'

const server = new StellarSdk.Server(process.env.REACT_APP_STELLAR_SERVER_URL)

type StellarBalanceResponse = {
  balances: {
    asset_type: string
    balance: string
  }[]
}

const connectStellarWallet = (setStellarAccount: Dispatch<string>): void => {
  getPublicKey()
    .then(publicKey => {
      setStellarAccount(publicKey)
    })
    .catch(error => toast.error(error.message))
}

const getBalanceAccount = async (publicKey: string): Promise<string> => {
  return server
    .accounts()
    .accountId(publicKey)
    .call()
    .then((accountResult: StellarBalanceResponse) => {
      const ethBalance = accountResult.balances.find(
        balance => balance.asset_type === 'native'
      )?.balance
      return Promise.resolve(ethBalance || 0)
    })
    .catch(() => {
      return Promise.reject('Unable to get Stellar Wallet balance!')
    })
}

const createPaymentTransaction = async (
  publicKey: string,
  publicKeyDestination: string,
  amount: string,
  ethereumKey: string
): Promise<string> => {
  try {
    const account = await server.loadAccount(publicKey)
    const fee = await server.fetchBaseFee()

    const keyArray = zeroPad(arrayify(ethereumKey), 32)

    const transaction = new StellarSdk.TransactionBuilder(account, {
      fee,
      networkPassphrase: StellarSdk.Networks.TESTNET,
    })
      .addOperation(
        StellarSdk.Operation.payment({
          destination: publicKeyDestination,
          asset: StellarSdk.Asset.native(),
          amount: amount,
        })
      )
      .addMemo(new StellarSdk.Memo.hash(Buffer.from(keyArray)))
      .setTimeout(0)
      .build()

    const xdr = transaction.toEnvelope().toXDR('base64')
    return Promise.resolve(xdr)
  } catch (error) {
    return Promise.reject(error)
  }
}

const signStellarTransaction = async (xdr: string): Promise<string> => {
  try {
    const signedTransaction = await signTransaction(xdr, {
      networkPassphrase: process.env.REACT_APP_STELLAR_NETWORK_PASSPHRASE
    })

    const transactionToSubmit = StellarSdk.TransactionBuilder.fromXDR(
      signedTransaction,
      process.env.REACT_APP_STELLAR_SERVER_URL
    )

    const transactionResult = await server.submitTransaction(
      transactionToSubmit
    )

    return Promise.resolve(transactionResult.hash)
  } catch (error) {
    return Promise.reject(error)
  }
}

const signMultipleStellarTransactions = async (
  xdrList: WithdrawResult[]
): Promise<string> => {
  try {
    const mainTransaction = StellarSdk.TransactionBuilder.fromXDR(
      xdrList[0].xdr,
      process.env.REACT_APP_STELLAR_SERVER_URL
    )

    for (let i = 1; i < Math.trunc(validatorUrls.length / 2) + 1; i++) {
      const transaction = StellarSdk.TransactionBuilder.fromXDR(
        xdrList[i].xdr,
        process.env.REACT_APP_STELLAR_SERVER_URL
      )
      mainTransaction.addDecoratedSignature(...transaction.signatures)
    }

    return Promise.resolve(mainTransaction.toXDR())
  } catch (error) {
    return Promise.reject(error)
  }
}

export {
  connectStellarWallet,
  createPaymentTransaction,
  signStellarTransaction,
  signMultipleStellarTransactions,
  getBalanceAccount,
}
