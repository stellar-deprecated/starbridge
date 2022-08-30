import React, { useState } from 'react'

import classNames from 'classnames'

import {
  CurrencyLabel,
  IInputProps,
} from 'components/molecules/labeled-input/type'

import Eth from 'app/core/resources/eth-icon.svg'
import Weth from 'app/core/resources/weth-icon.svg'

import { Label } from '../../atoms/typography/label'
import styles from './styles.module.scss'

export enum InputLabel {
  sending = 'Sending',
  receive = 'Receive',
}

export interface ILabeledInputProps extends IInputProps {
  htmlType?: string
  currency: CurrencyLabel
  label: InputLabel
}

export interface ICurrencyProps {
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

    const getCurrencyData = (): ICurrencyProps => {
      const data = {
        [CurrencyLabel.eth]: { label: CurrencyLabel.eth, iconPath: Eth },
        [CurrencyLabel.weth]: { label: CurrencyLabel.weth, iconPath: Weth },
      }
      return data[currency]
    }

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
                  <Label text={`Bal: 1.42 ${getCurrencyData().label}`} />
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
                src={getCurrencyData().iconPath}
                className={styles.icon}
                alt={currency}
              />
              <Label
                text={getCurrencyData().label}
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
