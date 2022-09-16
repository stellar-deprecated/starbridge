import React from 'react'
import { Route } from 'react-router-dom'

import { IModuleRouteProps } from '../types'

const createModuleRoutes = ({
  routePrefix = '',
  routes,
}: IModuleRouteProps): React.ReactNode[] => {
  return routes.map(({ component, path }) => {
    const routeProps = { path: `${routePrefix}${path}` }
    return <Route {...routeProps} element={React.createElement(component)} />
  })
}

export default createModuleRoutes
