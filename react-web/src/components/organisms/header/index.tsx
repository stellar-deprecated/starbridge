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
  const { sendingAccount, receivingAccount, logoutSending, logoutReceiving } =
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
          iconLeft={<img src={Eth} alt="Eth" />}
          variant={ButtonVariant.tertiary}
          disabled={!sendingAccount}
          onClick={logoutSending}
        >
          <Typography
            variant={TypographyVariant.label}
            text={handleButtonText(sendingAccount)}
            className={styles.labelButton}
          />
        </Button>

        <Button
          iconLeft={<img src={Weth} alt="Weth" />}
          className={styles.loginButton}
          variant={ButtonVariant.tertiary}
          disabled={!receivingAccount}
          onClick={logoutReceiving}
        >
          <Typography
            variant={TypographyVariant.label}
            text={handleButtonText(receivingAccount)}
            className={styles.labelButton}
          />
        </Button>
      </div>
    </header>
  )
}

export { Header }
