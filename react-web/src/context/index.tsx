import { createContext, Dispatch, useState, useContext } from 'react'

const useLocalStorage = (
  storageKey: string
): [string, Dispatch<string | undefined>] => {
  const [storedValue, setStoredValue] = useState(
    window.localStorage.getItem(storageKey) || ''
  )

  const setValue = (value: string | undefined): void => {
    const valueToStore = value || ''
    window.localStorage.setItem(storageKey, valueToStore)
    setStoredValue(valueToStore)
  }

  return [storedValue, setValue]
}

export type AuthProviderProps = {
  stellarAccount: string
  setStellarAccount: Dispatch<string>
  ethereumAccount: string
  setEthereumAccount: Dispatch<string>
  concordiumAccount: string
  setConcordiumAccount: Dispatch<string>
  ethereumProvider?: string
  logoutStellar: () => void
  logoutEthereum: () => void
  logoutConcordium: () => void
}

export const AuthContext = createContext<AuthProviderProps>({
  stellarAccount: '',
  setStellarAccount: () => {
    return ''
  },
  ethereumAccount: '',
  setEthereumAccount: () => {
    return ''
  },
  concordiumAccount: '',
  setConcordiumAccount: () => {
    return ''
  },
  ethereumProvider: undefined,
  logoutStellar: () => undefined,
  logoutEthereum: () => undefined,
  logoutConcordium: () => undefined,
})

type AuthContextProviderProps = {
  children: React.ReactNode
}

export const AuthContextProvider = ({
  children,
}: AuthContextProviderProps): JSX.Element => {
  const ethereumProvider =
    localStorage.getItem('WEB3_CONNECT_CACHED_PROVIDER') || undefined
  const [stellarAccount, setStellarAccount] = useLocalStorage('stellarAccount')
  const [ethereumAccount, setEthereumAccount] =
    useLocalStorage('ethereumAccount')
  const [concordiumAccount, setConcordiumAccount] =
    useLocalStorage('concordiumAccount')

  const logoutStellar = (): void => {
    setStellarAccount(undefined)
  }

  const logoutEthereum = (): void => {
    setEthereumAccount(undefined)
    window.localStorage.setItem('walletconnect', '')
  }

  const logoutConcordium = (): void => {
    setConcordiumAccount(undefined)
  }

  return (
    <AuthContext.Provider
      value={{
        stellarAccount: stellarAccount,
        setStellarAccount: setStellarAccount,
        ethereumAccount: ethereumAccount,
        concordiumAccount: concordiumAccount,
        ethereumProvider: ethereumProvider,
        setEthereumAccount: setEthereumAccount,
        setConcordiumAccount: setConcordiumAccount,
        logoutStellar: logoutStellar,
        logoutEthereum: logoutEthereum,
        logoutConcordium: logoutConcordium,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export const useAuthContext = (): AuthProviderProps => useContext(AuthContext)
