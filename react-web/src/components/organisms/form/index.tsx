import React, { useEffect } from 'react'
import { useForm } from 'react-hook-form'

import { yupResolver } from '@hookform/resolvers/yup'

import {
  LabeledInput,
  ILabeledInputProps,
} from 'components/molecules/labeled-input'
import { IInputProps } from 'components/molecules/labeled-input/type'

import { renderFormElement } from './helpers'
import { IFormInputProps, IFormProps } from './types'

const FORM_ELEMENT_MAP = {
  InputText: LabeledInput as React.FC<
    Omit<IFormInputProps & ILabeledInputProps, 'ref'>
  >,
}

const isFormComponent = (
  component: React.ReactElement<IInputProps>
): string | boolean =>
  Object.values(FORM_ELEMENT_MAP).some(
    el => el === (component as React.ReactElement).type
  ) && component.props.name

const FormComponent = ({
  children,
  onSubmit,
  validationSchema,
  onError,
}: IFormProps): JSX.Element => {
  const {
    register,
    handleSubmit,
    formState,
    formState: { errors, isSubmitted },
  } = useForm({
    resolver: validationSchema ? yupResolver(validationSchema) : undefined,
  })
  useEffect(() => {
    onError?.(formState.errors)
  }, [formState, onError])

  const renderFormElements = (
    elements: React.ReactElement | React.ReactNode
  ): React.ReactElement | React.ReactNode =>
    React.Children.toArray(elements).map(child => {
      if (React.isValidElement(child)) {
        if (isFormComponent(child)) {
          const formActions = { errors, register, isSubmitted }
          return renderFormElement(child, formActions)
        }
      }
      return child
    })

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {renderFormElements(children)}
    </form>
  )
}

const Form = Object.assign(FormComponent, FORM_ELEMENT_MAP)
export { Form }
