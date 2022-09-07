import React from 'react'
import {
  DeepMap,
  FieldError,
  FieldValues,
  UseFormRegister,
} from 'react-hook-form'

import { Status } from 'components/enums'
import { LabeledInput } from 'components/molecules/labeled-input'

import { UpdateFormEvent } from './types'

const composeEvents =
  (
    eventA: (e: UpdateFormEvent) => void,
    eventB?: (e: UpdateFormEvent) => void
  ) =>
  (...args: [UpdateFormEvent]): void => {
    eventB && eventB(...args)
    return eventA(...args)
  }

export const cloneElement = (
  child: React.ReactElement,
  register: UseFormRegister<FieldValues>,
  status: Status
): React.ReactElement => {
  const { onChange, onBlur, ...registerProps } = register(child.props.name)
  const composeOnChange = composeEvents(onChange, child.props?.onChange)
  const composeOnBlur = composeEvents(onBlur, child.props?.onBlur)
  const { labeledInputClassName, ...props } = child.props
  const elementToClone = { ...child, props }
  return React.cloneElement(elementToClone, {
    onChange: composeOnChange,
    onBlur: composeOnBlur,
    status,
    ...registerProps,
  })
}

interface IFormActions {
  register: UseFormRegister<FieldValues>
  errors: DeepMap<FieldValues, FieldError>
  isSubmitted: boolean
}

export const renderFormElement = (
  child: React.ReactElement,
  { errors, register, isSubmitted }: IFormActions
): React.ReactElement => {
  const name = child.props.name
  const error = errors[name]?.message
  const status =
    error !== undefined
      ? Status.error
      : isSubmitted
      ? Status.success
      : Status.default

  cloneElement(child, register, status)

  return (
    <LabeledInput
      key={child.props.currency}
      label={child.props.label}
      className={child.props.labeledInputClassName}
      name={child.props.name}
      currency={child.props.currency}
      placeholder="--"
    />
  )
}
