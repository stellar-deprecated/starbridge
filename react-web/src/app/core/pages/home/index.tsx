import { useState, useEffect } from 'react'
import { toast } from 'react-toastify'

import { useAuthContext } from 'context'
import Web3 from 'web3'
import Web3utils from 'web3-utils'

import { TransactionStep } from 'components/enums'
import { HomeTemplate } from 'components/templates/home'
import { Currency } from 'components/types/currency'

import {
  connectEthereumWallet,
  depositEthereumTransaction,
  withdrawEthereumTransaction,
} from 'interfaces/ethereum'
import { deposit, withdraw, WithdrawResult } from 'interfaces/http'
import {
  connectStellarWallet,
  createPaymentTransaction,
  signStellarTransaction,
  signMultipleStellarTransactions,
  getBalanceAccount,
} from 'interfaces/stellar'

const Home = (): JSX.Element => {
  const [isLoading, setIsLoading] = useState(false)
  const [isModalLoading, setIsModalLoading] = useState(false)
  const [transactionStep, setTransactionStep] = useState(
    TransactionStep.deposit
  )
  const [balanceStellarAccount, setBalanceStellarAccount] = useState('')
  const [balanceEthereumAccount, setBalanceEthereumAccount] = useState('')
  const [xdrPaymentTransaction, setXdrPaymentTransaction] = useState('')
  const [xdrWithdrawTransaction, setXdrWithdrawTransaction] = useState<
    WithdrawResult[]
  >([])
  const [transactionHash, setTransactionHash] = useState('')
  const [transactionLogIndex, setTransactionLogIndex] = useState('')
  const [currentTransactionDetails, setCurrentTransactionDetails] = useState('')

  const {
    stellarAccount,
    setStellarAccount,
    ethereumAccount,
    setEthereumAccount,
  } = useAuthContext()

  useEffect(() => {
    getCurrentStellarBalance(stellarAccount)
    getCurrentEthereumBalance(ethereumAccount)
  }, [ethereumAccount, stellarAccount])

  const resetPage = (): void => {
    setTransactionStep(TransactionStep.deposit)
    setXdrPaymentTransaction('')
    setXdrWithdrawTransaction([])
    setTransactionHash('')
    setTransactionLogIndex('')
  }

  const getCurrentStellarBalance = (publicKey: string): void => {
    if (!publicKey) {
      setBalanceStellarAccount('')
      return
    }

    getBalanceAccount(publicKey).then(balance => {
      setBalanceStellarAccount(parseFloat(balance).toFixed(7))
    })
  }

  const getCurrentEthereumBalance = async (
    publicKey: string
  ): Promise<void> => {
    if (!publicKey) {
      setBalanceEthereumAccount('')
      return
    }

    await window.ethereum.enable()
    const web3 = new Web3(window.ethereum)

    web3.eth.getBalance(publicKey).then(balance => {
      setBalanceEthereumAccount(
        parseFloat(Web3utils.fromWei(balance)).toFixed(7)
      )
    })
  }

  const connectWallet = async (currentCurrency: Currency): Promise<void> => {
    currentCurrency === Currency.WETH
      ? connectStellarWallet(setStellarAccount)
      : connectEthereumWallet(setEthereumAccount)
  }

  const handleSendingButtonClick = async (
    currencyFrom: Currency
  ): Promise<void> => {
    connectWallet(currencyFrom)
  }

  const handleReceivingButtonClick = (currencyTo: Currency): void => {
    connectWallet(currencyTo)
  }

  const createStellarPaymentTrasaction = (value: string): void => {
    createPaymentTransaction(
      stellarAccount,
      process.env.REACT_APP_STELLAR_BRIDGE_ACCOUNT || '',
      value,
      ethereumAccount
    )
      .then(xdr => {
        setXdrPaymentTransaction(xdr)
        setCurrentTransactionDetails(xdr)
        setTransactionStep(TransactionStep.signDeposit)
      })
      .finally(() => setIsLoading(false))
  }

  const createStellarWithdrawTransaction = async (): Promise<void> => {
    withdraw(Currency.WETH, transactionHash)
      .then(results => {
        withdrawEthereumTransaction(results, ethereumAccount)
          .then(() => {
            resetPage()
            toast.success(
              'The transfer to your Ethereum Wallet was successful!'
            )
          })
          .catch(error => {
            toast.error(error.message)
          })
          .finally(() => setIsLoading(false))
      })
      .catch(error => {
        toast.error(error.message)
        setIsLoading(false)
      })
  }

  const createEthereumWithdrawTransaction = (): void => {
    withdraw(Currency.ETH, transactionHash, transactionLogIndex)
      .then(results => {
        setXdrWithdrawTransaction(results)
        setTransactionStep(TransactionStep.signWithdraw)
      })
      .catch(error => {
        toast.error(error?.response?.data.detail)
      })
      .finally(() => setIsLoading(false))
  }

  const handleSubmit = async (
    value: string,
    currencyFlow: Currency
  ): Promise<void> => {
    setIsLoading(true)

    // eslint-disable-next-line no-console
    console.log("transactionStep before deposit", transactionStep)

    if (transactionStep === TransactionStep.deposit) {
      setIsLoading(true)
      if (currencyFlow === Currency.WETH) {
        createStellarPaymentTrasaction(value)
      } else {
        depositEthereumTransaction(stellarAccount, ethereumAccount, value)
          .then(result => {
            setTransactionHash(result.transactionHash)
            setTransactionLogIndex(result.events.Deposit.logIndex)

            deposit(stellarAccount, result.transactionHash, result.events.Deposit.logIndex)
              .then(() => setTransactionStep(TransactionStep.withdraw))
              .finally(() => setIsLoading(false))
          })
          .catch(error => {
            setIsLoading(false)
            toast.error(error.message)
          })
      }
    }

    // eslint-disable-next-line no-console
    console.log("transactionStep after deposit", transactionStep)

    if (transactionStep === TransactionStep.withdraw) {
      currencyFlow === Currency.WETH
        ? createStellarWithdrawTransaction()
        : createEthereumWithdrawTransaction()
    }
  }

  const handleDepositSignTransaction = (): void => {
    setIsModalLoading(true)
    signStellarTransaction(xdrPaymentTransaction)
      .then(transactionHash => {
        setTransactionHash(transactionHash)
        setTransactionStep(TransactionStep.withdraw)

        // deposit(stellarAccount, transactionHash, null)
        //   .then(() => {
        //     setTransactionStep(TransactionStep.withdraw)
        //   })
        //   .catch(error => {
        //     toast.error(error.message)
        //   })
      })
      .finally(() => setIsModalLoading(false))
  }

  const handleWithdrawSignTransaction = (): void => {
    setIsModalLoading(true)

    signMultipleStellarTransactions(xdrWithdrawTransaction)
      .then(xdr => {
        setCurrentTransactionDetails(xdr)
        signStellarTransaction(xdr)
          .then(() => {
            resetPage()
            toast.success('The transfer to your Stellar Wallet was successful!')
          })
          .catch(error => {
            toast.error(error.message)
          })
          .finally(() => setIsModalLoading(false))
      })
      .catch(error => {
        toast.error(error.message)
        setIsModalLoading(false)
      })
  }

  const handleCancelClick = (): void => {
    setTransactionStep(
      transactionStep === TransactionStep.signWithdraw
        ? TransactionStep.withdraw
        : TransactionStep.deposit
    )
  }

  return (
    <HomeTemplate
      isLoading={isLoading}
      isModalLoading={isModalLoading}
      transactionStep={transactionStep}
      transactionDetails={currentTransactionDetails}
      onSubmit={handleSubmit}
      balanceStellarAccount={balanceStellarAccount}
      balanceEthereumAccount={balanceEthereumAccount}
      onSendingButtonClick={handleSendingButtonClick}
      onReceivingButtonClick={handleReceivingButtonClick}
      onDepositSignTransaction={handleDepositSignTransaction}
      onWithdrawSignTransaction={handleWithdrawSignTransaction}
      onCancelClick={handleCancelClick}
    />
  )
}

export default Home
