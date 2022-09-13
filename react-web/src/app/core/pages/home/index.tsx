import { useForm } from 'react-hook-form'

import { yupResolver } from '@hookform/resolvers/yup'
import * as yup from 'yup'

import { HomeTemplate } from 'components/templates/home'

const Home = (): JSX.Element => {
  const validationSchema = yup.object().shape({
    amountReceived: yup.number().min(1).required('This field is required.'),
    amountSent: yup.number().min(1).required('This field is required.'),
  })

  const defaultValues: { amountReceived?: number; amountSent?: number } = {
    amountReceived: 0,
    amountSent: 0,
  }

  const { handleSubmit } = useForm({
    resolver: yupResolver(validationSchema),
    defaultValues,
  })

  const onSubmit = (): void => {
    console.log('submit')
  }
  return <HomeTemplate handleSubmit={handleSubmit(onSubmit)} />
}

export default Home
