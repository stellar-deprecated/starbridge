import React, { ChangeEvent, FocusEvent } from 'react'

import { AnyObjectSchema } from 'yup'

export interface IFormProps {
  children: React.ReactNode | React.ReactNode[]

  onSubmit: (
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    values: Record<string | number, any>,
    event: React.BaseSyntheticEvent | undefined
  ) => Promise<void> | void

  validationSchema?: AnyObjectSchema

  initialValues?: Record<string | number, unknown>

  onError?: (errors: Record<string | number, unknown>) => void
}

export interface IFormInputProps {
  name: string

  ref?: React.Ref<HTMLInputElement>

  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void

  onBlur?: (e: React.FocusEvent<HTMLInputElement>) => void

  label?: string

  labeledInputClassName?: string
}

export type UpdateFormEvent =
  | ChangeEvent<HTMLInputElement>
  | FocusEvent<HTMLInputElement>
