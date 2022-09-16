import { useState } from 'react'
import { useForm } from 'react-hook-form'

import { yupResolver } from '@hookform/resolvers/yup'
import * as yup from 'yup'

import { HomeTemplate } from 'components/templates/home'

import { openWalletConnector } from 'interfaces/wallet-connect'

const Home = (): JSX.Element => {
  //TODO: remove these eslint disable
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [sendingWalletAccount, setSendingWalletAccount] = useState('')
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [receivingWalletAccount, setReceivingWalletAccount] = useState('')

  const validationSchema = yup.object().shape({
    amountReceived: yup.number().min(1).required('This field is required.'),
    amountSent: yup.number().min(1).required('This field is required.'),
  })

  const defaultValues: { amountReceived?: number; amountSent?: number } = {
    amountReceived: 0,
    amountSent: 0,
  }

  const { handleSubmit } = useForm({
    resolver: yupResolver(validationSchema),
    defaultValues,
  })

  const handleSendingButtonClick = (): void => {
    openWalletConnector().then(setSendingWalletAccount)
  }

  const handleReceivingButtonClick = (): void => {
    openWalletConnector().then(setReceivingWalletAccount)
  }

  const onSubmit = (): void => {
    //TODO: add submit logic here
  }
  return (
    <HomeTemplate
      handleSubmit={handleSubmit(onSubmit)}
      onSendingButtonClick={handleSendingButtonClick}
      onReceivingButtonClick={handleReceivingButtonClick}
    />
  )
}

export default Home
