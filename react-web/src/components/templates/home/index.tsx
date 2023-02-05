import { useEffect, useState, useCallback } from 'react'

import { useAuthContext } from 'context'

import {
  Button,
  ButtonVariant,
  ButtonSize,
  Typography,
  TypographyVariant,
} from 'components/atoms'
import { TransactionStep } from 'components/enums'
import { ICurrencyProps, InputLabel } from 'components/molecules'
import { SignTransactionModal } from 'components/organisms'
import { WalletInput } from 'components/organisms/wallet-input'
import { Currency, CurrencyLabel } from 'components/types/currency'

import SwitchIcon from 'app/core/resources/switch.svg'

import styles from './styles.module.scss'

export interface IHomeTemplateProps {
  isLoading?: boolean
  isModalLoading?: boolean
  transactionStep?: TransactionStep
  balanceStellarAccount?: string
  balanceEthereumAccount?: string
  transactionDetails?: string
  onSubmit: (value: string, currencyFlow: Currency) => void
  onSendingButtonClick?: (currencyFrom: Currency) => void
  onReceivingButtonClick?: (currencyTo: Currency) => void
  onDepositSignTransaction?: () => void
  onWithdrawSignTransaction?: () => void
  onCancelClick?: () => void
}

const HomeTemplate = ({
  isLoading = false,
  isModalLoading = false,
  transactionStep = TransactionStep.deposit,
  balanceStellarAccount = '',
  balanceEthereumAccount = '',
  transactionDetails,
  onSubmit,
  onSendingButtonClick,
  onReceivingButtonClick,
  onDepositSignTransaction,
  onWithdrawSignTransaction,
  onCancelClick,
}: IHomeTemplateProps): JSX.Element => {
  const { stellarAccount, ethereumAccount } = useAuthContext()

  const [isButtonEnabled, setIsButtonEnabled] = useState(false)
  const [inputSent, setInputSent] = useState('')
  const [receiveValue, setReceiveValue] = useState('')
  const [currencyFrom, setCurrencyFrom] = useState(Currency.WETH)
  const [currencyTo, setCurrencyTo] = useState(Currency.ETH)
  const [isOpenModal, setIsOpenModal] = useState(false)
  const [errorMessage, setErrorMessage] = useState('')

  const isCurrentStep = useCallback(
    (currentStep: string[]): boolean => {
      return currentStep.includes(transactionStep)
    },
    [transactionStep]
  )

  const isDepositOrSignDepositStep = isCurrentStep([
    TransactionStep.deposit,
    TransactionStep.signDeposit,
  ])

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

  const handleErrorInput = (value: string, currency: Currency): void => {
    const amount = parseFloat(value) || 0

    if (amount <= 0) {
      setErrorMessage('This field is required')
    } else if (
      (currency === Currency.ETH &&
        amount > parseFloat(balanceEthereumAccount)) ||
      (currency === Currency.WETH && amount > parseFloat(balanceStellarAccount))
    ) {
      setErrorMessage('Amount exceeds wallet balance')
    } else {
      setErrorMessage('')
    }
  }

  const onInputSentChange = (evt: React.FormEvent<HTMLInputElement>): void => {
    const input = evt.target as HTMLInputElement

    handleErrorInput(input.value, currencyFrom)

    setInputSent(input.value)
    setReceiveValue(input.value)
  }

  useEffect(() => {
    setIsButtonEnabled(
      !!inputSent && !errorMessage && !!stellarAccount && !!ethereumAccount
    )
  }, [errorMessage, inputSent, stellarAccount, ethereumAccount])

  useEffect(() => {
    setIsOpenModal(
      isCurrentStep([TransactionStep.signDeposit, TransactionStep.signWithdraw])
    )
  }, [setIsOpenModal, isCurrentStep])

  const changeCurrency = (): void => {
    setCurrencyFrom(prev => {
      const newCurrencyFrom =
        prev === Currency.ETH ? Currency.WETH : Currency.ETH
      handleErrorInput(inputSent, newCurrencyFrom)
      return newCurrencyFrom
    })
    setCurrencyTo(prev =>
      prev === Currency.ETH ? Currency.WETH : Currency.ETH
    )
  }

  const handleSendingButtonClick = (): void => {
    onSendingButtonClick && onSendingButtonClick(currencyFrom)
  }

  const handleReceivingButtonClick = (): void => {
    onReceivingButtonClick && onReceivingButtonClick(currencyTo)
  }

  const handleSubmit = (): void => {
    onSubmit(receiveValue, currencyFrom)
  }

  const handleCancelModal = (): void => {
    setIsOpenModal(false)
    onCancelClick && onCancelClick()
  }

  return (
    <main className={styles.main}>
      <div className={styles.bgMask}></div>
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
              accountConnected={
                currencyFrom === Currency.WETH
                  ? stellarAccount
                  : ethereumAccount
              }
              balanceAccount={
                currencyFrom === Currency.WETH
                  ? balanceStellarAccount
                  : balanceEthereumAccount
              }
              onChange={onInputSentChange}
              name={InputLabel.sending}
              onClick={handleSendingButtonClick}
              alreadySubmittedDeposit={!isDepositOrSignDepositStep}
              errorMessage={errorMessage}
            />
          </div>
          <div className={styles.formRow}>
            <WalletInput
              currency={currencyPropsConverter[currencyTo]}
              accountConnected={
                currencyTo === Currency.WETH ? stellarAccount : ethereumAccount
              }
              name={InputLabel.receive}
              disabled
              placeholder={receiveValue ? receiveValue : '--'}
              onClick={handleReceivingButtonClick}
            />
          </div>
          <Button
            isLoading={isLoading}
            variant={ButtonVariant.primary}
            fullWidth
            disabled={!isButtonEnabled}
            onClick={handleSubmit}
          >
            {isDepositOrSignDepositStep ? 'Send Transfer' : 'Withdraw'}
          </Button>
        </div>
      </div>
      <SignTransactionModal
        isOpen={isOpenModal}
        isLoading={isModalLoading}
        setModalOpen={setIsOpenModal}
        title={`${
          isDepositOrSignDepositStep ? 'Deposit' : 'Withdraw'
        } Transaction`}
        platform={Currency.WETH}
        transactionDetails={transactionDetails}
        onSignTransactionClick={
          isDepositOrSignDepositStep
            ? onDepositSignTransaction
            : onWithdrawSignTransaction
        }
        onCancelClick={handleCancelModal}
      />
    </main>
  )
}

export { HomeTemplate }
