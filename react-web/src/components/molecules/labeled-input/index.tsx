import React, { useState } from 'react'

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
      ...restProps
    },
    ref
  ): JSX.Element => {
    const [hasBalanceInfo, setHasBalanceInfo] = useState(false)

    const renderCurrencyInfo = (
      evt: React.ChangeEvent<HTMLInputElement>
    ): void => {
      setHasBalanceInfo(!!evt.currentTarget.value)
      return onChange?.(evt)
    }

    return (
      <div
        className={classNames(
          styles.inputContainer,
          !isSender && styles.receiveContainer
        )}
      >
        <div className={styles.inputRow}>
          <Typography
            variant={TypographyVariant.label}
            text={label}
            className={styles.mainLabel}
          />
          <div>
            {hasBalanceInfo && label === InputLabel.sending && (
              <>
                <Typography
                  variant={TypographyVariant.label}
                  text="Set Max"
                  className={styles.balanceLabel}
                />
                <Typography
                  variant={TypographyVariant.label}
                  text={`Bal: 1.42 ${currency.initials}`}
                />
              </>
            )}
          </div>
        </div>
        <div className={classNames(styles.inputFooter, className)}>
          <input
            id={id ?? name}
            className={classNames(styles.input, className)}
            onChange={renderCurrencyInfo}
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
