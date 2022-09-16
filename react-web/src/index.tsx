import React from 'react'
import ReactDOM from 'react-dom'

import { Buffer } from 'buffer'

import { Header } from 'components/organisms/header'

import App from 'app/core/App'

import './index.css'

window.Buffer = window.Buffer || Buffer

ReactDOM.render(
  <React.StrictMode>
    <Header
      labelWalletButton={'Not Connected'}
      labelLoginButton={'Not Connected'}
    />
    <App />
  </React.StrictMode>,
  document.getElementById('root')
)
