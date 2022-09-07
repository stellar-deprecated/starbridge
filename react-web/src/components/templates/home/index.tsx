import {
  Button,
  ButtonVariant,
  ButtonSize,
  Typography,
  TypographyVariant,
} from 'components/atoms'
import { InputLabel, LabeledInput } from 'components/molecules'
import { CurrencyLabel } from 'components/types/currency'

import SwitchIcon from 'app/core/resources/switch.svg'

import styles from './styles.module.scss'

export interface IHomeTemplateProps {
  transactionTitle?: string
}

interface IWalletInputProps {
  currency: CurrencyLabel
  isSender?: boolean
}

const WalletInput = ({
  isSender,
  currency,
}: IWalletInputProps): JSX.Element => {
  const settings = {
    name: {
      [CurrencyLabel.weth]: 'stellar wallet input',
      [CurrencyLabel.eth]: 'ethereum wallet input',
    },
    label: isSender ? InputLabel.sending : InputLabel.receive,
    button: {
      [CurrencyLabel.weth]: 'Connect Stellar Wallet',
      [CurrencyLabel.eth]: 'Connect Ethereum Wallet',
    },
  }

  return (
    <div className={styles.inputContainer}>
      <div className={styles.content}>
        <Typography
          className={styles.label}
          variant={TypographyVariant.p}
          text={isSender ? 'From:' : 'To:'}
        />
        <Button variant={ButtonVariant.secondary} size={ButtonSize.small}>
          {settings.button[currency]}
        </Button>
      </div>
      <LabeledInput
        name={settings.name[currency]}
        label={settings.label}
        currency={currency}
      />
    </div>
  )
}

const HomeTemplate = ({
  transactionTitle = 'Ethereum -> Stellar',
}: IHomeTemplateProps): JSX.Element => {
  return (
    <main className={styles.main}>
      <div className={styles.container}>
        <div className={styles.titleContainer}>
          <Typography
            className={styles.title}
            variant={TypographyVariant.h3}
            text={transactionTitle}
          />
          <Button
            className={styles.button}
            variant={ButtonVariant.primary}
            size={ButtonSize.small}
            iconLeft={<img src={SwitchIcon} alt="Switch Icon" />}
          >
            Switch
          </Button>
        </div>
        <WalletInput isSender currency={CurrencyLabel.eth} />
        <WalletInput currency={CurrencyLabel.weth} />
        <Button variant={ButtonVariant.primary} fullWidth>
          Send Transfer
        </Button>
      </div>
    </main>
  )
}

export { HomeTemplate }
