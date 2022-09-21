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

import styles from './styles.module.scss'

interface IWalletInputProps extends IInputProps {
  currency: ICurrencyProps
  isSender?: boolean
  accountConnected?: string
  onClick?: () => void
}

export const WalletInput = ({
  isSender,
  currency,
  name,
  placeholder,
  accountConnected,
  onChange,
  onClick,
}: IWalletInputProps): JSX.Element => {
  return (
    <div className={styles.inputContainer}>
      <div className={styles.content}>
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
      <LabeledInput
        name={currency.initials}
        label={isSender ? InputLabel.sending : InputLabel.receive}
        currency={currency}
        onChange={onChange}
        isSender={isSender}
        disabled={name === InputLabel.receive}
        placeholder={placeholder ?? '--'}
      />
    </div>
  )
}
