import axios from 'axios'

import {Currency} from 'components/types/currency'
import {chain} from "lodash";

export const validatorUrls = [
  process.env.REACT_APP_STARBRIDGE_VALIDATOR_URL_1,
  // process.env.REACT_APP_STARBRIDGE_VALIDATOR_URL_2,
  // process.env.REACT_APP_STARBRIDGE_VALIDATOR_URL_3,
]

export type WithdrawResult = {
  xdr: string
  address: string
  deposit_id: string
  expiration: string
  signature: string
  amount: number
  token: string
}

const deposit = async (
  account: string,
  transactionHash: string,
  transactionLogIndex: string
): Promise<string> => {
  return new Promise<string>((resolve, reject) => {
    const form = new FormData()
    form.append('hash', transactionHash)
    form.append('stellar_address', account)
    form.append('log_index', transactionLogIndex)

    const promises = validatorUrls.map(url =>
      axios.post(`${url}/ethereum/deposit`, form)
    )

    Promise.all(promises)
      .then(result => {
        return resolve(result[0].data)
      })
      .catch(error => {
        return reject(error)
      })
  })
}

const withdraw = async (
  currency: Currency,
  chain: Currency,
  transactionHash: string,
  transactionIndex = ''
): Promise<WithdrawResult[]> => {
  console.log(currency)
  console.log(chain)
  const isFromStellar = currency === Currency.WETH
  const isEthereumChain = chain === Currency.ETH
  const isConcordiumChain = chain === Currency.WCCD
  console.log(isFromStellar)
  console.log(isEthereumChain)
  console.log(isConcordiumChain)
  const form = new FormData()
  form.append('transaction_hash', transactionHash)

  if (!isFromStellar && isEthereumChain) {
    form.append('log_index', transactionIndex)
  }

  const promises = validatorUrls.map(url => {
    return new Promise<WithdrawResult>((resolve, reject) => {
      const postWithdraw = async (): Promise<void> => {
        try {
          let fullUrl = ""
          if (isFromStellar) {
            if (isEthereumChain) {
              fullUrl = `${url}/stellar/withdraw/ethereum`
            } else if (isConcordiumChain) {
              fullUrl = `${url}/stellar/withdraw/concordium`
            }
          } else {
            if (isEthereumChain) {
              fullUrl = `${url}/ethereum/withdraw/stellar`
            } else if (isConcordiumChain) {
              fullUrl = `${url}/concordium/withdraw/stellar`
            }
          }
          const response = await axios.post(
              fullUrl,
            form
          )

          switch (response.status) {
            case 202:
              break
            case 200:
              resolve(isFromStellar ? response.data : { xdr: response.data })
              return
            default:
              reject(response)
              return
          }

          setTimeout(() => {
            postWithdraw()
          }, 2000)
        } catch (error) {
          reject(error)
        }
      }

      postWithdraw()
    })
  })

  return Promise.all(promises)
}

export { deposit, withdraw }
