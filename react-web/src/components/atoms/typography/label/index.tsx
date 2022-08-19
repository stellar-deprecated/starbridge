import classNames from 'classnames'

import { FontSize } from 'components/enums/font-size'
import { Status } from 'components/enums/status'

import { ITypographyDefaultProps } from '..'
import styles from './style.module.scss'

export interface ILabelProps extends ITypographyDefaultProps {
  /**
   * Status style for label
   */
  status?: Status
  /**
   * Label for element
   */
  htmlFor?: string
  /**
   * The label font size
   */
  fontSize?: FontSize
}

const Label = ({
  text = '',
  status = Status.default,
  fontSize = FontSize.normal,
  className,
}: ILabelProps): JSX.Element => {
  return (
    <span
      className={classNames(
        styles.label,
        className,
        styles[fontSize],
        styles[`${status}-color`]
      )}
    >
      {text}
    </span>
  )
}

export { Label }
