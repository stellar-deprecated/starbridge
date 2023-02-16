import axios from 'axios'

import {Currency} from 'components/types/currency'

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
  chain: Currency,
  account: string,
  transactionHash: string,
  transactionLogIndex = ""
): Promise<string> => {
  return new Promise<string>((resolve, reject) => {
    const form = new FormData()
    form.append('hash', transactionHash)
    form.append('stellar_address', account)
    if (transactionLogIndex !== "") {
      form.append('log_index', transactionLogIndex)
    }

    let promises
    if (chain === Currency.ETH){
      promises = validatorUrls.map(url =>
        axios.post(`${url}/ethereum/deposit`, form)
      )
    } else {
      // if (chain === Currency.WCCD)
      promises = validatorUrls.map(url =>
          axios.post(`${url}/concordium/deposit`, form)
      )
    }

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
  destinationAddress: string,
  transactionIndex = '',
): Promise<WithdrawResult[]> => {
  const isFromStellar = currency === Currency.WETH
  const isEthereumChain = chain === Currency.ETH
  const isConcordiumChain = chain === Currency.WCCD
  const form = new FormData()
  form.append('transaction_hash', transactionHash)

  if (!isFromStellar && isEthereumChain) {
    form.append('log_index', transactionIndex)
  }
  if (chain === Currency.WCCD){
    form.append('destination', destinationAddress)
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
