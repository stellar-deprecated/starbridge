import React from 'react'
import { FieldValues, useForm } from 'react-hook-form'

import { InputText } from 'components/atoms/input'

import { CoreRouter } from 'app/core/routes'

const App = (): JSX.Element => {
  const { register, handleSubmit, watch } = useForm({})

  return (
    <InputText
      {...register('eventInformation')}
      placeholder="Event information"
      // className={styles.input}
    />
  )
}

export default App
