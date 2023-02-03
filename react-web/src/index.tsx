import React from 'react'
import ReactDOM from 'react-dom'
import { ToastContainer } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'

import { Buffer } from 'buffer'
import { AuthContextProvider } from 'context'

import { Header } from 'components/organisms/header'

import App from 'app/core/App'

import './index.css'

window.Buffer = window.Buffer || Buffer

ReactDOM.render(
  <React.StrictMode>
    <AuthContextProvider>
      <Header />
      <App />
      <ToastContainer />
    </AuthContextProvider>
  </React.StrictMode>,
  document.getElementById('root')
)
