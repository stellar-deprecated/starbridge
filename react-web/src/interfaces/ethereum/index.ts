import { Dispatch } from 'react'
import { toast } from 'react-toastify'

import { BigNumber } from '@ethersproject/bignumber'
import StellarSdk from 'stellar-sdk'
import Web3 from 'web3'
import Web3utils from 'web3-utils'
import Web3Modal from 'web3modal'

import BridgeContractBuild from 'interfaces/ethereum/Bridge.json'
import { validatorUrls, WithdrawResult } from 'interfaces/http'

type EthereumDepositContractResult = {
  transactionHash: string
  events: {
    Deposit: {
      logIndex: string
    }
  }
}

const web3Modal = new Web3Modal({
  cacheProvider: true,
})

type Web3ProviderProps = {
  selectedAddress: string
}

const connectEthereumWallet = (setEthereumAccount: Dispatch<string>): void => {
  const currentProvider = web3Modal.connect()

  currentProvider
    .then((result: Web3ProviderProps) => {
      setEthereumAccount(result.selectedAddress)
    })
    .catch(error => {
      toast.error(error.message)
    })
}

const depositEthereumTransaction = async (
  stellarAccount: string,
  ethereumAccount: string,
  value: string
): Promise<EthereumDepositContractResult> => {
  await window.ethereum.enable()
  const web3 = new Web3(window.ethereum)

  const bridgeContract = new web3.eth.Contract(
    BridgeContractBuild.abi as Web3utils.AbiItem[],
    process.env.REACT_APP_ETHEREUM_BRIDGE_ACCOUNT
  )

  const stellarAccountDecoded =
    StellarSdk.StrKey.decodeEd25519PublicKey(stellarAccount)

  return bridgeContract.methods
    .depositETH(BigNumber.from(stellarAccountDecoded))
    .send({
      from: ethereumAccount,
      value: BigNumber.from(Web3utils.toWei(value)),
    })
}

const withdrawEthereumTransaction = async (
  withdrawResult: WithdrawResult[],
  ethereumAccount: string
): Promise<void> => {
  await window.ethereum.enable()
  const web3 = new Web3(window.ethereum)

  const withdrawERC20Request = {
    id: `0x${withdrawResult[0].deposit_id}`,
    expiration: BigNumber.from(withdrawResult[0].expiration),
    recipient: ethereumAccount,
    amount: withdrawResult[0].amount,
  }

  const bridgeContract = new web3.eth.Contract(
    BridgeContractBuild.abi as Web3utils.AbiItem[],
    process.env.REACT_APP_ETHEREUM_BRIDGE_ACCOUNT
  )

  const indexes: number[] = []
  const signatures: Buffer[] = []

  for (let i = 0; i < validatorUrls.length; i++) {
    const addressSigner: string = await bridgeContract.methods.signers(i).call()
    const currentSignature =
      withdrawResult.find(result => result.address === addressSigner)
        ?.signature || ''
    indexes.push(i)
    signatures.push(Buffer.from(currentSignature, 'hex'))
  }

  return bridgeContract.methods
    .withdrawETH(withdrawERC20Request, signatures, indexes)
    .send({ from: ethereumAccount })
}

export {
  connectEthereumWallet,
  depositEthereumTransaction,
  withdrawEthereumTransaction,
}
