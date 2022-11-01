import React from 'react'

import classNames from 'classnames'

import { WalletIcon, Typography, TypographyVariant } from 'components/atoms'
import { IInputProps } from 'components/molecules/labeled-input/type'
import { Currency, CurrencyLabel } from 'components/types/currency'

import styles from './styles.module.scss'

export enum InputLabel {
  sending = 'Sending',
  receive = 'Receive',
}

export interface ILabeledInputProps extends IInputProps {
  htmlType?: string
  currency: ICurrencyProps
  label: InputLabel
  isSender?: boolean
  balanceAccount?: string
  errorMessage?: string
}

export interface ICurrencyProps {
  initials: CurrencyLabel
  label: Currency
}

const LabeledInput = React.forwardRef<HTMLInputElement, ILabeledInputProps>(
  (
    {
      name,
      onChange,
      htmlType = 'number',
      className,
      id,
      currency,
      label,
      isSender,
      placeholder,
      balanceAccount,
      errorMessage,
      ...restProps
    },
    ref
  ): JSX.Element => {
    return (
      <div
        className={classNames(
          styles.inputContainer,
          !isSender && styles.receiveContainer,
          errorMessage && styles.error
        )}
      >
        {errorMessage && (
          <Typography
            className={styles.errorLabel}
            variant={TypographyVariant.label}
            text={errorMessage}
          />
        )}
        <div className={styles.inputRow}>
          <Typography
            variant={TypographyVariant.label}
            text={label}
            className={styles.mainLabel}
          />
          {balanceAccount && label === InputLabel.sending && (
            <div>
              <Typography
                variant={TypographyVariant.label}
                text="Set Max"
                className={styles.balanceLabel}
              />
              <Typography
                variant={TypographyVariant.label}
                text={`Bal: ${balanceAccount} ${currency.initials}`}
              />
            </div>
          )}
        </div>
        <div className={classNames(styles.inputFooter, className)}>
          <input
            id={id ?? name}
            className={classNames(styles.input, className)}
            onChange={onChange}
            type={htmlType}
            name={name}
            placeholder={placeholder ?? '--'}
            {...restProps}
            ref={ref}
          />

          <div className={styles.currencyContainer}>
            <WalletIcon currency={currency.label} />
            <Typography
              variant={TypographyVariant.label}
              text={currency.initials}
              className={styles.currencyLabel}
            />
          </div>
        </div>
      </div>
    )
  }
)

export { LabeledInput }
