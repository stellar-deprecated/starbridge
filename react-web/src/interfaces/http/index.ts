import axios from 'axios'

const deposit = (account: string, transactionHash: string): void => {
  const validatorUrls = [
    'https://starbridge1.prototypes.kube001.services.stellar-ops.com',
    'https://starbridge2.prototypes.kube001.services.stellar-ops.com',
    'https://starbridge3.prototypes.kube001.services.stellar-ops.com',
  ]

  const form = new FormData()
  form.append('hash', transactionHash)
  form.append('stellar_address', account)

  const promises = validatorUrls.map(url => axios.post(`${url}/deposit`, form))

  Promise.all(promises)
    .then(results => {
      console.log('results', results)
    })
    .catch(error => {
      console.log('error', error)
    })
}

export { deposit }
