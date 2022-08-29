import { Route } from 'react-router'
import { Link } from 'react-router-dom'

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
}

const Header = ({
  className,
  handleOnWalletPress,
  handleOnLoginPress,
}: IHeaderProps): JSX.Element => {
  return (
    <header className={classNames(styles.header, className)}>
      <nav>
        {/* <Link to="/"> */}
        <Heading
          level={TypographyHeadingLevel.h3}
          text="Starbridge"
          className={styles.title}
        />
        {/* </Link> */}
      </nav>
      <div className={styles.containerButton}>
        <Button
          iconLeft={<img src={Eth} alt="Eth" />}
          variant={ButtonVariant.tertiary}
          onClick={(): void => {
            handleOnWalletPress
          }}
        >
          <Label text="0xb79...9268 " className={styles.labelButton} />
        </Button>

        <Button
          iconLeft={<img src={Weth} alt="Weth" />}
          className={styles.loginButton}
          variant={ButtonVariant.tertiary}
          onClick={(): void => {
            handleOnLoginPress
          }}
        >
          <Label text="Not Connected " className={styles.labelButton} />
        </Button>
      </div>
    </header>
  )
}

export { Header }
