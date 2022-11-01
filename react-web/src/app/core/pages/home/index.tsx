import { useState } from 'react'

import WalletConnect from '@walletconnect/client'
import { useAuthContext } from 'context'
import { convertStringToHex, sanitizeHex } from 'utils'

import { HomeTemplate } from 'components/templates/home'

import { deposit } from 'interfaces/http'
import { createPaymentTransaction } from 'interfaces/stellar'
import { openWalletConnector } from 'interfaces/wallet-connect'

const Home = (): JSX.Element => {
  const [connector, setConnector] = useState<WalletConnect>()
  const { sendingAccount, setSendingAccount } = useAuthContext()

  const handleSendingButtonClick = (): void => {
    openWalletConnector().then(connector => {
      setConnector(connector)

      connector.on('connect', (error, payload) => {
        //TODO: do more tests with different wallets and their returns
        const { accounts } = payload.params[0]
        console.log('payload', payload)
        setSendingAccount(accounts[0])
      })
    })
  }

  const handleReceivingButtonClick = (): void => {
    // openWalletConnector().then(setReceivingAccount)
  }

  const handleSubmit = async (): Promise<void> => {
    const from = sendingAccount || ''
    const to = sendingAccount || ''
    const value = sanitizeHex(convertStringToHex(0))
    const data = '0x'

    const tx = {
      from,
      to,
      value,
      data,
    }

    console.log('tx', tx)
    try {
      console.log('start')
      // send transaction
      if (!connector) {
        console.log('not connected')
        return
      }

      console.log('start send')
      const result = await connector.sendTransaction(tx)
      console.log('result - sendTransaction', result)
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <HomeTemplate
      onSubmit={handleSubmit}
      onSendingButtonClick={handleSendingButtonClick}
      onReceivingButtonClick={handleReceivingButtonClick}
    />
  )
}

export default Home
