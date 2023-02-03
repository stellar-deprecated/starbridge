import { Button as SButton, Loader } from '@stellar/design-system'
import classNames from 'classnames'

import styles from './styles.module.scss'

export enum ButtonSize {
  default = 'default',
  small = 'small',
}

export enum ButtonVariant {
  primary = 'primary',
  secondary = 'secondary',
  tertiary = 'tertiary',
  ghost = 'ghost',
}

export interface IButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  className?: string
  iconLeft?: React.ReactNode
  iconRight?: React.ReactNode
  variant?: ButtonVariant
  size?: ButtonSize
  isLoading?: boolean
  fullWidth?: boolean
  disabled?: boolean
  onClick?: React.MouseEventHandler<HTMLButtonElement>
  children: string | React.ReactNode
}

const Button = (props: IButtonProps): JSX.Element => {
  const {
    variant = ButtonVariant.primary,
    size = ButtonSize.default,
    fullWidth,
    className,
    isLoading,
    disabled,
    children,
    ...rest
  } = props
  const fullWidthStyle = fullWidth && styles.fullWidth

  return (
    <SButton
      className={classNames(
        styles.button,
        styles[variant],
        styles[size],
        fullWidthStyle,
        className
      )}
      disabled={disabled || isLoading}
      {...rest}
    >
      {isLoading ? <Loader size="20px" /> : children}
    </SButton>
  )
}

export { Button }
