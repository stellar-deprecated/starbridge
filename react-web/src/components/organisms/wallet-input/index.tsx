import { formatWalletAccount } from 'utils'

import {
  Button,
  ButtonVariant,
  ButtonSize,
  Typography,
  TypographyVariant,
  WalletIcon,
} from 'components/atoms'
import { ICurrencyProps, InputLabel, LabeledInput } from 'components/molecules'
import { IInputProps } from 'components/molecules/labeled-input/type'

import SuccessIcon from 'app/core/resources/success.svg'

import styles from './styles.module.scss'

interface IWalletInputProps extends IInputProps {
  currency: ICurrencyProps
  isSender?: boolean
  accountConnected?: string
  balanceAccount?: string
  alreadySubmittedDeposit?: boolean
  errorMessage?: string
  onClick?: () => void
}

export const WalletInput = ({
  isSender,
  currency,
  name,
  placeholder,
  accountConnected,
  balanceAccount,
  alreadySubmittedDeposit,
  errorMessage,
  onChange,
  onClick,
}: IWalletInputProps): JSX.Element => {
  return (
    <div className={styles.inputContainer}>
      <div className={styles.content}>
        <div className={styles.insideContent}>
          <Typography
            className={styles.label}
            variant={TypographyVariant.p}
            text={isSender ? 'From:' : 'To:'}
          />
          {accountConnected ? (
            <div className={styles.accountContainer}>
              <WalletIcon currency={currency.label} />
              <Typography
                className={styles.account}
                variant={TypographyVariant.p}
                text={formatWalletAccount(accountConnected)}
              />
            </div>
          ) : (
            <Button
              variant={ButtonVariant.secondary}
              size={ButtonSize.small}
              onClick={onClick}
            >
              {`Connect ${currency.label} Wallet`}
            </Button>
          )}
        </div>
        {isSender && alreadySubmittedDeposit && (
          <div className={styles.insideContent}>
            <Typography
              className={styles.submittedTxn}
              variant={TypographyVariant.p}
              text="Submitted Txn"
            />
            <img src={SuccessIcon} alt="Success Icon" />
          </div>
        )}
      </div>
      <LabeledInput
        id={name}
        name={currency.initials}
        label={isSender ? InputLabel.sending : InputLabel.receive}
        currency={currency}
        onChange={onChange}
        isSender={isSender}
        disabled={name === InputLabel.receive}
        placeholder={placeholder ?? '--'}
        balanceAccount={balanceAccount}
        errorMessage={errorMessage}
      />
    </div>
  )
}
