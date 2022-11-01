import BigNumber from 'bignumber.js'

const formatWalletAccount = (account: string): string => {
  const accountLength = account.length
  return account
    ? `${account.substring(0, 5)}...${account.substring(
        accountLength - 4,
        accountLength
      )}`
    : account
}

const sanitizeHex = (hex: string): string => {
  hex = hex.substring(0, 2) === '0x' ? hex.substring(2) : hex
  if (hex === '') {
    return ''
  }
  hex = hex.length % 2 !== 0 ? '0' + hex : hex
  return '0x' + hex
}

const convertStringToHex = (value: string | number): string => {
  return new BigNumber(`${value}`).toString(16)
}

export { formatWalletAccount, sanitizeHex, convertStringToHex }
