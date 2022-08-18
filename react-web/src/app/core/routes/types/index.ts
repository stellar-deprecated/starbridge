import React from 'react'

export type AppRoute = {
  path: string
  isPrivate?: boolean
  component: React.FunctionComponent | React.ComponentClass
}

export interface IModuleRouteProps {
  routePrefix?: string
  routes: AppRoute[]
}
