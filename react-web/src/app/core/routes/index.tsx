import React from 'react'
import { BrowserRouter as Router, Routes } from 'react-router-dom'

import createModuleRoutes from './module-routes'
import { coreRoutes } from './routes'

const CoreRouter = (): JSX.Element => {
  const routes = [
    createModuleRoutes({
      routePrefix: '',
      routes: coreRoutes,
    }),
  ]

  return (
    <Router>
      <Routes>{routes.map(route => route)}</Routes>
    </Router>
  )
}

export { CoreRouter }
