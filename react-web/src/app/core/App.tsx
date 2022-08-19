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
  </>
)

export default App
