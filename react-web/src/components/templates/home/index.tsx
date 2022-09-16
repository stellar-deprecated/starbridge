import { FormEvent, useEffect, useRef, useState } from 'react'

import {
  Button,
  ButtonVariant,
  ButtonSize,
  Typography,
  TypographyVariant,
} from 'components/atoms'
import { ICurrencyProps } from 'components/molecules'
import { WalletInput } from 'components/organisms/wallet-input'
import { Currency, CurrencyLabel } from 'components/types/currency'

import Eth from 'app/core/resources/eth.svg'
import SwitchIcon from 'app/core/resources/switch.svg'
import Weth from 'app/core/resources/weth.svg'

import styles from './styles.module.scss'

export interface IHomeTemplateProps {
  transactionTitle?: string
  handleSubmit: (evt: FormEvent<HTMLFormElement>) => Promise<void>
  onSendingButtonClick?: () => void
  onReceivingButtonClick?: () => void
}

const HomeTemplate = ({
  handleSubmit,
  onSendingButtonClick,
  onReceivingButtonClick,
}: IHomeTemplateProps): JSX.Element => {
  const [isButtonEnabled, setIsButtonEnabled] = useState(false)
  const sendingRef = useRef<HTMLInputElement>(null)
  const receivingRef = useRef<HTMLInputElement>(null)
  const [inputSent, setInputSent] = useState('')
  const [inputReceived, setInputReceived] = useState('')
  const [currencyFrom, setCurrencyFrom] = useState(Currency.ETH)
  const [currencyTo, setCurrencyTo] = useState(Currency.WETH)

  const currencyPropsConverter: Record<Currency, ICurrencyProps> = {
    [Currency.WETH]: {
      initials: CurrencyLabel.weth,
      label: Currency.WETH,
      iconPath: Weth,
    },
    [Currency.ETH]: {
      initials: CurrencyLabel.eth,
      label: Currency.ETH,
      iconPath: Eth,
    },
  }

  const onInputReceivedChange = (
    evt: React.FormEvent<HTMLInputElement>
  ): void => {
    const input = evt.target as HTMLInputElement
    setInputReceived(input.value)
  }

  const onInputSentChange = (evt: React.FormEvent<HTMLInputElement>): void => {
    const input = evt.target as HTMLInputElement

    setInputSent(input.value)
  }

  useEffect(() => {
    setIsButtonEnabled(inputSent > '0' && inputReceived > '0')
  }, [inputSent, inputReceived])

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
          <form data-testid="form" onSubmit={handleSubmit}>
            <div className={styles.formRow}>
              <WalletInput
                isSender
                currency={currencyPropsConverter[currencyFrom]}
                onChange={onInputSentChange}
                name={'sending'}
                onClick={onSendingButtonClick}
                ref={sendingRef}
              />
            </div>
            <div className={styles.formRow}>
              <WalletInput
                currency={currencyPropsConverter[currencyTo]}
                onChange={onInputReceivedChange}
                name={'receiving'}
                onClick={onReceivingButtonClick}
                ref={receivingRef}
              />
            </div>
            <Button
              variant={ButtonVariant.primary}
              fullWidth
              disabled={!isButtonEnabled}
            >
              Send Transfer
            </Button>
          </form>
        </div>
      </div>
    </main>
  )
}

export { HomeTemplate }
