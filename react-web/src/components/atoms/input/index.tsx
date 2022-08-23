import React, { SVGProps, FC } from 'react'

import classNames from 'classnames'

import { Status } from 'components/enums'

import Error from 'app/core/resources/svgs/error.svg'
import Success from 'app/core/resources/svgs/success.svg'

import styles from './styles.module.scss'
import { InputProps } from './type'

export interface IInputTextProps extends InputProps {
  htmlType?: string
  customIcon?: JSX.Element
}

interface IInputIconProps {
  status: Status
  customIcon?: JSX.Element
}

const InputIcon = ({
  status = Status.default,
  customIcon,
}: IInputIconProps): JSX.Element | null => {
  if (customIcon) {
    return customIcon
  }
  return null
  // const iconMap: { [key: string]: React.FC<SVGProps<SVGSVGElement>> } = {
  //   // [Status.success]: Success,
  //   // [Status.error]: Error,
  // }
  // const Icon = iconMap[status]
  // if (!Icon) {
  //   return null
  // }
  // return <Icon className={styles.icon} />
}

const InputText = React.forwardRef<HTMLInputElement, IInputTextProps>(
  (
    {
      name,
      onChange,
      onBlur,
      htmlType = 'text',
      disabled = false,
      className,
      id,
      status = Status.default,
      customIcon,
      ...restProps
    },
    ref
  ): JSX.Element => {
    return (
      <div className={classNames(styles.inputContainer)}>
        <input
          id={id ?? name}
          className={classNames(
            styles.input,
            styles[`${status}Input`],
            {
              [styles.disabled]: disabled,
            },
            className
          )}
          onChange={onChange}
          onBlur={onBlur}
          type={htmlType}
          name={name}
          {...restProps}
          ref={ref}
        />
        <InputIcon status={status} customIcon={customIcon} />
      </div>
    )
  }
)

export { InputText }
