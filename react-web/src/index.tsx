import React from 'react'
import ReactDOM from 'react-dom'

import '@stellar/design-system/build/styles.min.css'

import { Header } from 'components/organisms/header'

import App from 'app/core/App'

import './index.css'

ReactDOM.render(
  <React.StrictMode>
    <Header />
    <App />
  </React.StrictMode>,
  document.getElementById('root')
)
