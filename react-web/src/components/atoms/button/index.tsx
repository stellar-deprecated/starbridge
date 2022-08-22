import { Button as SButton } from '@stellar/design-system'
import classNames from 'classnames'

import styles from './styles.module.scss'

export enum ButtonVariant {
  primary = 'primary',
  secondary = 'secondary',
  tertiary = 'tertiary',
}

export interface IButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  className?: string
  iconLeft?: React.ReactNode
  iconRight?: React.ReactNode
  variant?: ButtonVariant
  isLoading?: boolean
  fullWidth?: boolean
  disabled?: boolean
  children: string | React.ReactNode
}

const Button = (props: IButtonProps): JSX.Element => {
  const {
    variant = ButtonVariant.primary,
    fullWidth,
    className,
    ...rest
  } = props
  const fullWidthStyle = fullWidth && styles.fullWidth

  return (
    <SButton
      className={classNames(
        styles.button,
        styles[variant],
        fullWidthStyle,
        className
      )}
      {...rest}
    />
  )
}

export { Button }
