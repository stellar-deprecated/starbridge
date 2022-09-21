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
  sendingAccount?: string
  setSendingAccount: Dispatch<string | undefined>
  receivingAccount?: string
  setReceivingAccount: Dispatch<string | undefined>
  logoutSending: () => void
  logoutReceiving: () => void
}

export const AuthContext = createContext<AuthProviderProps>({
  sendingAccount: undefined,
  setSendingAccount: () => {
    return undefined
  },
  receivingAccount: undefined,
  setReceivingAccount: () => {
    return undefined
  },
  logoutSending: () => undefined,
  logoutReceiving: () => undefined,
})

type AuthContextProviderProps = {
  children: React.ReactNode
}

export const AuthContextProvider = ({
  children,
}: AuthContextProviderProps): JSX.Element => {
  const [sendingAccount, setSendingAccount] = useLocalStorage('sendingAccount')
  const [receivingAccount, setReceivingAccount] =
    useLocalStorage('receivingAccount')

  const logoutSending = (): void => {
    setSendingAccount(undefined)
  }

  const logoutReceiving = (): void => {
    setReceivingAccount(undefined)
  }

  return (
    <AuthContext.Provider
      value={{
        sendingAccount,
        setSendingAccount,
        receivingAccount,
        setReceivingAccount,
        logoutSending,
        logoutReceiving,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export const useAuthContext = (): AuthProviderProps => useContext(AuthContext)
