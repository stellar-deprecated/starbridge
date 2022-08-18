import { Button as SButton } from '@stellar/design-system'

export enum ButtonVariant {
  primary = 'primary',
  secondary = 'secondary',
  tertiary = 'tertiary',
}

export enum ButtonSize {
  default = 'default',
  small = 'small',
}

export interface IButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  iconLeft?: React.ReactNode
  iconRight?: React.ReactNode
  variant?: ButtonVariant
  isLoading?: boolean
  size?: ButtonSize
  fullWidth?: boolean
  children: string | React.ReactNode
}

const Button = (props: IButtonProps): JSX.Element => {
  return <SButton {...props} />
}

export { Button }
