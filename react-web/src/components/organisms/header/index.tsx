import { BrowserRouter as Router, Link } from 'react-router-dom'

import classNames from 'classnames'

import { Button, ButtonVariant } from 'components/atoms'
import {
  Heading,
  TypographyHeadingLevel,
} from 'components/atoms/typography/heading'
import { Label } from 'components/atoms/typography/label'

import Eth from 'app/core/resources/eth.svg'
import Weth from 'app/core/resources/weth.svg'

import styles from './style.module.scss'

interface IHeaderProps {
  className?: string
  handleOnWalletPress?: () => void
  handleOnLoginPress?: () => void
  labelWalletButton: string
  labelLoginButton: string
}

const Header = ({
  className,
  handleOnWalletPress,
  handleOnLoginPress,
  labelWalletButton,
  labelLoginButton,
}: IHeaderProps): JSX.Element => {
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
          disabled
          onClick={(): void => {
            handleOnWalletPress
          }}
        >
          <Label text={labelWalletButton} className={styles.labelButton} />
        </Button>

        <Button
          iconLeft={<img src={Weth} alt="Weth" />}
          className={styles.loginButton}
          variant={ButtonVariant.tertiary}
          disabled
          onClick={(): void => {
            handleOnLoginPress
          }}
        >
          <Label text={labelLoginButton} className={styles.labelButton} />
        </Button>
      </div>
    </header>
  )
}

export { Header }
