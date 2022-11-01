import { useState } from 'react'

import classNames from 'classnames'

import {
  Button,
  ButtonVariant,
  Typography,
  TypographyVariant,
} from 'components/atoms'
import { Modal, IModalProps } from 'components/organisms'

import { ReactComponent as StarbridgeIcon } from 'app/core/resources/starbridge.svg'

import styles from './styles.module.scss'

export interface ISignTransactionModalProps
  extends Pick<IModalProps, 'isOpen' | 'setModalOpen'> {
  title: string
  platform: string
  isLoading?: boolean
  transactionDetails?: string
  onSignTransactionClick?: () => void
  onCancelClick?: () => void
}

const SignTransactionModal = (
  props: ISignTransactionModalProps
): JSX.Element => {
  const {
    title,
    platform,
    isLoading,
    transactionDetails,
    onSignTransactionClick,
    onCancelClick,
    ...rest
  } = props

  const [copiedOpacityStyle, setCopiedOpacityStyle] = useState('')

  const handleCopyTransaction = (): void => {
    transactionDetails && navigator.clipboard.writeText(transactionDetails)
    setCopiedOpacityStyle(styles.opacityFull)

    setTimeout(() => {
      setCopiedOpacityStyle('')
    }, 3000)
  }

  return (
    <Modal className={styles.modal} {...rest}>
      <StarbridgeIcon />
      <Typography
        className={styles.title}
        variant={TypographyVariant.h1}
        text={title}
      />

      <div className={styles.platformContainer}>
        <div className={styles.symbol} />
        <Typography
          className={styles.platform}
          variant={TypographyVariant.h4}
          text={platform}
        />
      </div>

      <div
        className={styles.transactionContainer}
        onClick={handleCopyTransaction}
      >
        <Typography variant={TypographyVariant.p} text={transactionDetails} />
        <Typography
          className={classNames(styles.copiedLabel, copiedOpacityStyle)}
          variant={TypographyVariant.p}
          text="Copied"
        />
      </div>

      <div className={styles.buttonsContainer}>
        <Button
          isLoading={isLoading}
          variant={ButtonVariant.secondary}
          onClick={onSignTransactionClick}
          fullWidth
        >
          Sign Transaction
        </Button>
        <Button variant={ButtonVariant.ghost} onClick={onCancelClick} fullWidth>
          Cancel
        </Button>
      </div>
    </Modal>
  )
}

export { SignTransactionModal }
