import { useAuthContext } from 'context'

import { HomeTemplate } from 'components/templates/home'

import { openWalletConnector } from 'interfaces/wallet-connect'

const Home = (): JSX.Element => {
  const { setSendingAccount, setReceivingAccount } = useAuthContext()

  const handleSendingButtonClick = (): void => {
    openWalletConnector().then(setSendingAccount)
  }

  const handleReceivingButtonClick = (): void => {
    openWalletConnector().then(setReceivingAccount)
  }

  const handleSubmit = (): void => {
    //TODO: add submit logic here
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
