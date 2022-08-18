import Home from 'app/core/pages/home'

import { AppRoute } from './types'

export const coreRoutes: AppRoute[] = [
  { path: '/', component: Home },
  { path: '/private', component: Home, isPrivate: true },
]
