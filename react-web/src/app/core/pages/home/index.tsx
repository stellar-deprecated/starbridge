import React from 'react'

import logo from 'app/core/resources/logo.svg'

import styles from './styles.module.scss'

const Home = (): JSX.Element => {
  return (
    <main>
      <header className={styles.header}>
        <img src={logo} className={styles.logo} alt="logo" />
        <p>Welcome to CKL Boilerplate</p>
        <a href="https://reactjs.org" target="_blank" rel="noopener noreferrer">
          Learn React
        </a>
      </header>
    </main>
  )
}

export default Home
