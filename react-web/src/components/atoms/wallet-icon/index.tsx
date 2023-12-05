import { Currency } from 'components/types/currency'

import Eth from 'app/core/resources/eth.svg'
import Weth from 'app/core/resources/weth.svg'

const WalletIconPath = {
  [Currency.ETH]: Eth,
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
