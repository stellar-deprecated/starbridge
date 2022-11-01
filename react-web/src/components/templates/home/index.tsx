import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'

import { yupResolver } from '@hookform/resolvers/yup'
import { useAuthContext } from 'context'
import * as yup from 'yup'

import {
  Button,
  ButtonVariant,
  ButtonSize,
  Typography,
  TypographyVariant,
} from 'components/atoms'
import { ICurrencyProps, InputLabel } from 'components/molecules'
import { WalletInput } from 'components/organisms/wallet-input'
import { Currency, CurrencyLabel } from 'components/types/currency'

import SwitchIcon from 'app/core/resources/switch.svg'

import styles from './styles.module.scss'

export interface IHomeTemplateProps {
  transactionTitle?: string
  onSubmit: () => void
  onSendingButtonClick?: () => void
  onReceivingButtonClick?: () => void
}

const HomeTemplate = ({
  onSubmit,
  onSendingButtonClick,
  onReceivingButtonClick,
}: IHomeTemplateProps): JSX.Element => {
  const { sendingAccount, receivingAccount } = useAuthContext()

  const [isButtonEnabled, setIsButtonEnabled] = useState(false)
  const [inputSent, setInputSent] = useState('')
  const [receiveValue, setReceiveValue] = useState('')
  const [currencyFrom, setCurrencyFrom] = useState(Currency.ETH)
  const [currencyTo, setCurrencyTo] = useState(Currency.WETH)

  const currencyPropsConverter: Record<Currency, ICurrencyProps> = {
    [Currency.WETH]: {
      initials: CurrencyLabel.weth,
      label: Currency.WETH,
    },
    [Currency.ETH]: {
      initials: CurrencyLabel.eth,
      label: Currency.ETH,
    },
  }

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

  const onInputSentChange = (evt: React.FormEvent<HTMLInputElement>): void => {
    const input = evt.target as HTMLInputElement

    setInputSent(input.value)
    setReceiveValue(input.value)
  }

  useEffect(() => {
    setIsButtonEnabled(inputSent > '0')
  }, [inputSent])

  const changeCurrency = (): void => {
    setCurrencyFrom(prev =>
      prev === Currency.ETH ? Currency.WETH : Currency.ETH
    )
    setCurrencyTo(prev =>
      prev === Currency.ETH ? Currency.WETH : Currency.ETH
    )
  }

  return (
    <main className={styles.main}>
      <div className={styles.container}>
        <div className={styles.titleContainer}>
          <Typography
            className={styles.title}
            variant={TypographyVariant.h3}
            text={`${currencyFrom} -> ${currencyTo}`}
          />
          <Button
            className={styles.button}
            variant={ButtonVariant.primary}
            size={ButtonSize.small}
            iconLeft={<img src={SwitchIcon} alt="Switch Icon" />}
            onClick={changeCurrency}
          >
            Switch
          </Button>
        </div>
        <div className={styles.form}>
          <div className={styles.formRow}>
            <WalletInput
              isSender
              currency={currencyPropsConverter[currencyFrom]}
              accountConnected={sendingAccount}
              onChange={onInputSentChange}
              name={InputLabel.sending}
              onClick={onSendingButtonClick}
            />
          </div>
          <div className={styles.formRow}>
            <WalletInput
              currency={currencyPropsConverter[currencyTo]}
              accountConnected={receivingAccount}
              name={InputLabel.receive}
              disabled
              placeholder={receiveValue ? receiveValue : '--'}
              onClick={onReceivingButtonClick}
            />
          </div>
          <Button
            variant={ButtonVariant.primary}
            fullWidth
            disabled={!isButtonEnabled}
            onClick={onSubmit}
          >
            Send Transfer
          </Button>
        </div>
      </div>
    </main>
  )
}

export { HomeTemplate }
