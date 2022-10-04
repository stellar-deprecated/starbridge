import WalletConnect from '@walletconnect/client'
import QRCodeModal from '@walletconnect/qrcode-modal'

const openWalletConnector = async (): Promise<WalletConnect> => {
  return new Promise<WalletConnect>(resolve => {
    const walletConnector = new WalletConnect({
      bridge: 'https://bridge.walletconnect.org',
      qrcodeModal: QRCodeModal,
    })

    if (walletConnector.connected) {
      walletConnector.killSession().then(() => walletConnector.createSession())
    } else {
      walletConnector.createSession()
    }

    resolve(walletConnector)
  })
}

export { openWalletConnector }
