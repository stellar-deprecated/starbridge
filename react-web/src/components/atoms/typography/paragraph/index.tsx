import classNames from 'classnames'

import { ITypographyDefaultProps } from '..'
import styles from './styles.module.scss'

export enum TypographyWeight {
  normal = 'normal',
  semiBold = 'semiBold',
  medium = 'medium',
  light = 'light',
  bold = 'bold',
}

export interface IParagraphProps extends ITypographyDefaultProps {
  /**
   The font weight of the text
  */
  weight?: TypographyWeight
}

const Paragraph = ({
  text,
  weight = TypographyWeight.normal,
  className,
}: IParagraphProps): JSX.Element => (
  <p className={classNames(styles.paragraph, styles[weight], className)}>
    {text}
  </p>
)

export { Paragraph }
