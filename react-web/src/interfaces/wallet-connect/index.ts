import WalletConnect from '@walletconnect/client'
import QRCodeModal from '@walletconnect/qrcode-modal'

const openWalletConnector = async (): Promise<string> => {
  return new Promise<string>((resolve, reject) => {
    const walletConnector = new WalletConnect({
      bridge: 'https://bridge.walletconnect.org',
      qrcodeModal: QRCodeModal,
    })

    if (walletConnector.connected) {
      walletConnector.killSession().then(() => walletConnector.createSession())
    } else {
      walletConnector.createSession()
    }

    walletConnector.on('connect', (error, payload) => {
      if (error) {
        reject(error)
      }

      //TODO: do more tests with different wallets and their returns
      const { accounts } = payload.params[0]
      resolve(accounts[0])
    })
  })
}

export { openWalletConnector }
