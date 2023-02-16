import {useEffect, useState} from 'react'
import {toast} from 'react-toastify'

import {useAuthContext} from 'context'

import {TransactionStep} from 'components/enums'
import {HomeTemplate} from 'components/templates/home'
import {Currency} from 'components/types/currency'

import {connectEthereumWallet, depositEthereumTransaction, withdrawEthereumTransaction,} from 'interfaces/ethereum'
import {
  connectConcordiumWallet,
  depositConcordiumTransaction,
  withdrawConcordiumTransaction
} from 'interfaces/concordium'
import {deposit, withdraw, WithdrawResult} from 'interfaces/http'
import {
  connectStellarWallet,
  createPaymentTransaction,
  getBalanceAccount,
  signMultipleStellarTransactions,
  signStellarTransaction,
} from 'interfaces/stellar'
import axios from "axios";

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
    concordiumAccount,
    setEthereumAccount,
    setConcordiumAccount,
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

    // await window.ethereum.request({ method: 'eth_requestAccounts' })
    // const web3 = new Web3(window.ethereum)

    // We can use this for proposal user add our token WGBM in his wallet
    // const assetAdded = await window.ethereum.request({
    //       method: 'wallet_watchAsset',
    //       params: {
    //         type: 'ERC20',
    //         options: {
    //           address: process.env.REACT_APP_ETHEREUM_TOKEN_ACCOUNT,
    //           symbol: 'GBM',
    //           decimals: 7,
    //           // image: 'https://foo.io/token-image.svg',
    //         },
    //       },
    //     })
    // const stellarAssetContract = new web3.eth.Contract(
    //     StellarAssetContractBuild.abi as Web3utils.AbiItem[],
    //     process.env.REACT_APP_ETHEREUM_TOKEN_ACCOUNT
    // )

    // await stellarAssetContract.methods.balanceOf(publicKey).call().then(balance => {
    //   setBalanceEthereumAccount(
    //       parseFloat(Web3utils.fromWei((Number(balance)*10**11).toString())).toFixed(7)
    //   )
    // });
  }

  const connectWallet = async (currentCurrency: Currency): Promise<void> => {
    if (currentCurrency === Currency.WETH){
      connectStellarWallet(setStellarAccount)
    } else if (currentCurrency == Currency.ETH) {
      connectEthereumWallet(setEthereumAccount)
    } else if (currentCurrency == Currency.WCCD) {
      connectConcordiumWallet(setConcordiumAccount)
    }
  }

  const handleSendingButtonClick = async (
    currencyFrom: Currency
  ): Promise<void> => {
    connectWallet(currencyFrom)
  }

  const handleReceivingButtonClick = (currencyTo: Currency): void => {
    connectWallet(currencyTo)
  }

  const createStellarPaymentTrasaction = (value: string, chain: Currency): void => {
    const destinationAccount = chain === Currency.ETH ? ethereumAccount : concordiumAccount
    createPaymentTransaction(
      stellarAccount,
      process.env.REACT_APP_STELLAR_BRIDGE_ACCOUNT || '',
      value,
        destinationAccount
    )
      .then(xdr => {
        setXdrPaymentTransaction(xdr)
        setCurrentTransactionDetails(xdr)
        setTransactionStep(TransactionStep.signDeposit)
      })
      .finally(() => setIsLoading(false))
  }

  const createStellarWithdrawTransaction = async (chain: Currency): Promise<void> => {
    withdraw(Currency.WETH, chain, transactionHash, concordiumAccount)
      .then(results => {
        if (chain === Currency.ETH){
          withdrawEthereumTransaction(results, ethereumAccount)
              .then(() => {
                resetPage()
                toast.success(
                    'The transfer to your Polygon Wallet was successful!'
                )
              })
              .catch(error => {
                toast.error(error.message)
              })
              .finally(() => setIsLoading(false))
        } else if (chain === Currency.WCCD) {
          withdrawConcordiumTransaction(results, concordiumAccount)
              .then(() => {
                resetPage()
                toast.success(
                    'The transfer to your Concordium Wallet was successful!'
                )
              })
              .catch(error => {
                toast.error(error.message)
              })
              .finally(() => setIsLoading(false))
        }
      })
      .catch(error => {
        toast.error(error.message)
        setIsLoading(false)
      })
  }

  const createEthereumWithdrawTransaction = (chain: Currency): void => {
    withdraw(Currency.ETH, chain, transactionHash, ethereumAccount, transactionLogIndex)
      .then(results => {
        setCurrentTransactionDetails(results[0].xdr)
        setXdrWithdrawTransaction(results)
        setTransactionStep(TransactionStep.signWithdraw)
      })
      .catch(error => {
        toast.error(error?.response?.data.detail)
      })
      .finally(() => setIsLoading(false))
  }

  const createConcordiumWithdrawTransaction = (chain: Currency): void => {
    withdraw(Currency.WCCD, chain, transactionHash, stellarAccount)
      .then(results => {
          setCurrentTransactionDetails(results[0].xdr)
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
    currencyFlow: Currency,
    chainFlow: Currency,
  ): Promise<void> => {
    setIsLoading(true)

    if (transactionStep === TransactionStep.deposit) {
      setIsLoading(true)
      if (currencyFlow === Currency.WETH) {
        createStellarPaymentTrasaction(value, chainFlow)
      } else if (currencyFlow === Currency.WCCD) {
        depositConcordiumTransaction(stellarAccount, concordiumAccount, value)
            .then(txHash => {
                axios.post(
                    'http://localhost:8130/invokeContract/getDepositParams',
                    {hash: txHash},
                    {headers: {"Content-Type": "application/json", "Access-Control-Allow-Origin": "*"}}
                ).then(()=>{
                    setTransactionHash(txHash)
                    deposit(chainFlow, stellarAccount, txHash)
                        .then(() => setTransactionStep(TransactionStep.withdraw))
                        .finally(()=>setIsLoading(false))
                })
            })
      } else {
        depositEthereumTransaction(stellarAccount, ethereumAccount, value)
          .then(result => {
            setTransactionHash(result.transactionHash)
            setTransactionLogIndex(result.events.Deposit.logIndex)

            deposit(chainFlow, stellarAccount, result.transactionHash, result.events.Deposit.logIndex)
              .then(() => setTransactionStep(TransactionStep.withdraw))
              .finally(() => setIsLoading(false))
          })
          .catch(error => {
            setIsLoading(false)
            toast.error(error.message)
          })
      }
    }
    if (transactionStep === TransactionStep.withdraw) {
      if (currencyFlow === Currency.WETH){
        if (chainFlow === Currency.ETH) {
          createStellarWithdrawTransaction(chainFlow)
        } else if (chainFlow === Currency.WCCD) {
          createStellarWithdrawTransaction(chainFlow)
        }
      } else if (currencyFlow === Currency.ETH){
        if (chainFlow === Currency.ETH) {
          createEthereumWithdrawTransaction(chainFlow)
        }
      } else if (currencyFlow === Currency.WCCD){
        if (chainFlow === Currency.WCCD) {
          createConcordiumWithdrawTransaction(chainFlow)
        }
      }
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
