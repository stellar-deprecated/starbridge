import { Typography, TypographyVariant } from 'components/atoms'
import {
  Heading,
  TypographyHeadingLevel,
} from 'components/atoms/typography/heading'
import { Label } from 'components/atoms/typography/label'
import { Link } from 'components/atoms/typography/link'
import { Paragraph } from 'components/atoms/typography/paragraph'

import { CoreRouter } from 'app/core/routes'

const App = (): JSX.Element => (
  <>
    <CoreRouter />
    <Heading level={TypographyHeadingLevel.h1} text="Hsomething" />
    <Label text="Label" />
    <Link text="Link" />
    <Paragraph text="Paragraph" />
    <Typography variant={TypographyVariant.h1} text="h1" />
    <Typography variant={TypographyVariant.h2} text="h2" />
    <Typography variant={TypographyVariant.h3} text="h3" />
    <Typography variant={TypographyVariant.h4} text="h4" />
    <Typography variant={TypographyVariant.h5} text="h5" />
    <Typography variant={TypographyVariant.h6} text="h6" />
    <Typography variant={TypographyVariant.label} text="Label" />
    <Typography variant={TypographyVariant.link} text="Link" />
    <Typography variant={TypographyVariant.p} text="Paragraph" />
  </>
)

export default App
