import React, { ReactNode } from 'react'

import { Heading, TypographyHeadingLevel } from './heading'
import { Label, ILabelProps } from './label'
import { Link, ILinkProps } from './link'
import { Paragraph, IParagraphProps } from './paragraph'

export enum TypographyVariant {
  label = 'label',
  link = 'link',
  p = 'p',
  h1 = 'h1',
  h2 = 'h2',
  h3 = 'h3',
  h4 = 'h4',
  h5 = 'h5',
  h6 = 'h6',
}

export interface ITypographyDefaultProps {
  /**
   * Text for Typography
   */
  text: ReactNode | string
  /**
   * Classname to adds custom css
   * */
  className?: string
}

export interface ITypographyProps
  extends ITypographyDefaultProps,
    ILabelProps,
    ILinkProps,
    IParagraphProps {
  /**
   * Variant for Typography
   */
  variant: TypographyVariant
}

const TypographyVariantComponent = React.memo(
  ({ text, variant, ...props }: ITypographyProps) => {
    switch (variant) {
      case TypographyVariant.h1:
      case TypographyVariant.h2:
      case TypographyVariant.h3:
      case TypographyVariant.h4:
      case TypographyVariant.h5:
      case TypographyVariant.h6:
        return (
          <Heading
            text={text}
            level={
              TypographyHeadingLevel[
                variant as keyof typeof TypographyHeadingLevel
              ]
            }
            {...props}
          />
        )

      case TypographyVariant.label:
        return <Label text={text} {...props} />

      case TypographyVariant.link:
        return <Link text={text} {...props} />

      case TypographyVariant.p:
        return <Paragraph text={text} {...props} />

      default:
        return <span>{text}</span>
    }
  }
)

const Typography = ({
  text,
  variant,
  ...props
}: ITypographyProps): JSX.Element => (
  <TypographyVariantComponent text={text} variant={variant} {...props} />
)

export { Typography }
