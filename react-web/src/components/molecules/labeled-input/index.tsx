import React, { useState } from 'react'

import classNames from 'classnames'

import { IInputProps } from 'components/molecules/labeled-input/type'

import { Label } from '../../atoms/typography/label'
import styles from './styles.module.scss'

export enum InputLabel {
  sending = 'Sending',
  receive = 'Receive',
}

export interface ILabeledInputProps extends IInputProps {
  htmlType?: string
  currency: ICurrencyProps
  label: InputLabel
}

export interface ICurrencyProps {
  initials: string
  label: string
  iconPath: string
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
      ...restProps
    },
    ref
  ): JSX.Element => {
    const [hasBalanceInfo, setHasBalanceInfo] = useState(false)

    const renderCurrencyInfo = (
      evt: React.ChangeEvent<HTMLInputElement>
    ): void => {
      if (evt.currentTarget.value) {
        setHasBalanceInfo(true)
      } else {
        setHasBalanceInfo(false)
      }
      return onChange?.(evt)
    }

    return (
      <div className={styles.container}>
        <div className={classNames(styles.inputContainer)}>
          <div className={styles.inputRow}>
            <Label text={label} className={styles.mainLabel} />
            <div>
              {hasBalanceInfo && (
                <>
                  <Label text="Set Max" className={styles.balanceLabel} />
                  <Label text={`Bal: 1.42 ${currency.initials}`} />
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
              placeholder="--"
              {...restProps}
              ref={ref}
            />

            <div className={styles.currencyContainer}>
              <img
                src={currency.iconPath}
                className={styles.icon}
                alt={currency.label}
              />
              <Label
                text={currency.initials}
                className={styles.currencyLabel}
              />
            </div>
          </div>
        </div>
      </div>
    )
  }
)

export { LabeledInput }
