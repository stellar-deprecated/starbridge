import axios from 'axios'

import { Currency } from 'components/types/currency'

export const validatorUrls = [
  process.env.REACT_APP_STARBRIDGE_VALIDATOR_URL_1,
  process.env.REACT_APP_STARBRIDGE_VALIDATOR_URL_2,
  process.env.REACT_APP_STARBRIDGE_VALIDATOR_URL_3,
]

export type WithdrawResult = {
  xdr: string
  address: string
  deposit_id: string
  expiration: string
  signature: string
  amount: number
}

const deposit = async (
  account: string,
  transactionHash: string
): Promise<string> => {
  return new Promise<string>((resolve, reject) => {
    const form = new FormData()
    form.append('hash', transactionHash)
    form.append('stellar_address', account)

    const promises = validatorUrls.map(url =>
      axios.post(`${url}/deposit`, form)
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
  transactionHash: string,
  transactionIndex = ''
): Promise<WithdrawResult[]> => {
  const isFromStellar = currency === Currency.WETH
  const form = new FormData()
  form.append('transaction_hash', transactionHash)

  if (!isFromStellar) {
    form.append('log_index', transactionIndex)
  }

  const promises = validatorUrls.map(url => {
    return new Promise<WithdrawResult>((resolve, reject) => {
      const postWithdraw = async (): Promise<void> => {
        try {
          const response = await axios.post(
            `${url}/${
              isFromStellar
                ? 'stellar/withdraw/ethereum'
                : 'ethereum/withdraw/stellar'
            }`,
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
