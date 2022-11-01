const formatWalletAccount = (account: string): string => {
  const accountLength = account.length
  return account
    ? `${account.substring(0, 5)}...${account.substring(
        accountLength - 4,
        accountLength
      )}`
    : account
}

export { formatWalletAccount }
