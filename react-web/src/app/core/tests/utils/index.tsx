import { render, RenderOptions, RenderResult } from '@testing-library/react'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'

// If the application has providers, you can add them in the wrapper below
const ApplicationProviders: React.FC<React.PropsWithChildren> = ({
  children,
}) => <MemoryRouter>{children}</MemoryRouter>

const customRender = (
  ui: React.ReactElement,
  options?: Omit<RenderOptions, 'queries'>
): RenderResult => render(ui, { wrapper: ApplicationProviders, ...options })

export * from '@testing-library/react'
export { customRender as render }
