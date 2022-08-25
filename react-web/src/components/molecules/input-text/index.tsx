import React, { useState } from 'react'

import classNames from 'classnames'

import { IInputProps } from 'components/molecules/input-text/type'

import Eth from 'app/core/resources/eth-icon.svg'

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

export interface IInputTextProps extends IInputProps {
  htmlType?: string
  currency: CurrencyLabel
  label: InputLabel
}

const InputText = React.forwardRef<HTMLInputElement, IInputTextProps>(
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

    return (
      <div style={{ maxWidth: 600, margin: 32 }}>
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
          <div
            className={classNames(
              styles.inputFooter,

              className
            )}
          >
            <input
              id={id ?? name}
              className={classNames(styles.input, className)}
              onChange={(): void => setHasBalanceInfo(true)}
              type={htmlType}
              name={name}
              {...restProps}
              ref={ref}
            />

            <div className={styles.currencyContainer}>
              <img src={Eth} className={styles.icon} alt={currency} />
              <Label text={currency} className={styles.currencyLabel} />
            </div>
          </div>
        </div>
      </div>
    )
  }
)

export { InputText }
