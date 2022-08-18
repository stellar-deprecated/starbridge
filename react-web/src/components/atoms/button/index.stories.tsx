import { Story, Meta } from '@storybook/react'
import React from 'react'

import { ReactComponent as Logo } from 'app/core/resources/logo.svg'

import { Button, ButtonVariant, ButtonSize, IButtonProps } from '.'

export default {
  title: 'Atoms/Button',
  component: Button,
} as Meta

const Template: Story<IButtonProps> = args => {
  const { iconLeft, iconRight, ...rest } = args
  return (
    <Button
      {...rest}
      iconLeft={iconLeft && <Logo />}
      iconRight={iconRight && <Logo />}
    />
  )
}

export const Default = Template.bind({})

Default.argTypes = {
  children: { control: 'text', defaultValue: 'Click me' },
  isLoading: { control: 'text', defaultValue: 'Click me' },
  size: {
    control: 'select',
    options: ButtonSize,
    defaultValue: ButtonSize.default,
  },
  variant: {
    control: 'select',
    options: ButtonVariant,
    defaultValue: ButtonVariant.primary,
  },
  iconLeft: { control: 'boolean', defaultValue: false },
  iconRight: { control: 'boolean', defaultValue: false },
  fullWidth: { control: 'boolean', defaultValue: false },
}
