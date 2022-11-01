import { BrowserRouter as Router, Link } from 'react-router-dom'

import classNames from 'classnames'
import { useAuthContext } from 'context'
import { formatWalletAccount } from 'utils'

import {
  Button,
  ButtonVariant,
  Typography,
  TypographyVariant,
} from 'components/atoms'
import {
  Heading,
  TypographyHeadingLevel,
} from 'components/atoms/typography/heading'

import Eth from 'app/core/resources/eth.svg'
import Weth from 'app/core/resources/weth.svg'

import styles from './style.module.scss'

interface IHeaderProps {
  className?: string
}

const Header = ({ className }: IHeaderProps): JSX.Element => {
  const { stellarAccount, ethereumAccount, logoutStellar, logoutEthereum } =
    useAuthContext()

  const handleButtonText = (account: string | undefined): string => {
    return account ? formatWalletAccount(account) : 'Not Connected'
  }

  return (
    <header className={classNames(styles.header, className)}>
      <nav>
        <Router>
          <Link to="/">
            <Heading
              level={TypographyHeadingLevel.h3}
              text="Starbridge"
              className={styles.title}
            />
          </Link>
        </Router>
      </nav>
      <div className={styles.containerButton}>
        <Button
          iconLeft={<img src={Weth} alt="Weth" />}
          variant={ButtonVariant.tertiary}
          disabled={!stellarAccount}
          onClick={logoutStellar}
        >
          <Typography
            variant={TypographyVariant.label}
            text={handleButtonText(stellarAccount)}
            className={styles.labelButton}
          />
        </Button>

        <Button
          iconLeft={<img src={Eth} alt="Eth" />}
          className={styles.loginButton}
          variant={ButtonVariant.tertiary}
          disabled={!ethereumAccount}
          onClick={logoutEthereum}
        >
          <Typography
            variant={TypographyVariant.label}
            text={handleButtonText(ethereumAccount)}
            className={styles.labelButton}
          />
        </Button>
      </div>
    </header>
  )
}

export { Header }
