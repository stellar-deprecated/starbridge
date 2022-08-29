import React, { useState } from 'react'

import classNames from 'classnames'

import { IInputProps } from 'components/molecules/labeled-input/type'

import Eth from 'app/core/resources/eth-icon.svg'
import Weth from 'app/core/resources/weth-icon.svg'

import { Label } from '../../atoms/typography/label'
import styles from './styles.module.scss'

export enum InputLabel {
  sending = 'Sending',
  receive = 'Receive',
}

export enum CurrencyLabel {
  eth = 'ETH',
  weth = 'WETH',
}

export interface ILabeledInputProps extends IInputProps {
  htmlType?: string
  currency: CurrencyLabel
  label: InputLabel
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
      evt: React.FormEvent<HTMLInputElement>
    ): void => {
      evt.currentTarget.value
        ? setHasBalanceInfo(true)
        : setHasBalanceInfo(false)
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
                  <Label text={`Bal: 1.42 ${currency}`} />
                </>
              )}
            </div>
          </div>
          <div className={classNames(styles.inputFooter, className)}>
            <input
              id={id ?? name}
              className={classNames(styles.input, className)}
              onChange={(evt): void => renderCurrencyInfo(evt)}
              type={htmlType}
              name={name}
              {...restProps}
              ref={ref}
            />

            <div className={styles.currencyContainer}>
              <img
                src={currency === CurrencyLabel.eth ? Eth : Weth}
                className={styles.icon}
                alt={currency}
              />
              <Label text={currency} className={styles.currencyLabel} />
            </div>
          </div>
        </div>
      </div>
    )
  }
)

export { LabeledInput }
