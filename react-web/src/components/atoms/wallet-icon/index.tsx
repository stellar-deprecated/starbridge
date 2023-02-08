import { Currency } from 'components/types/currency'

import Weth from 'app/core/resources/gbm.svg'
import Matic from 'app/core/resources/matic.svg'

const WalletIconPath = {
  [Currency.ETH]: Matic,
  [Currency.WETH]: Weth,
}

export type WalletIconProps = {
  className?: string
  currency: Currency
}

const WalletIcon = ({ className, currency }: WalletIconProps): JSX.Element => {
  return (
    <img className={className} src={WalletIconPath[currency]} alt={currency} />
  )
}

export { WalletIcon }
