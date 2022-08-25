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
    transactionDetails,
    onSignTransactionClick,
    onCancelClick,
    ...rest
  } = props

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

      <div className={styles.transactionContainer}>
        <Typography variant={TypographyVariant.p} text={transactionDetails} />
      </div>

      <div className={styles.buttonsContainer}>
        <Button
          variant={ButtonVariant.secondary}
          onClick={onSignTransactionClick}
          fullWidth
        >
          Sign Transaciton
        </Button>
        <Button variant={ButtonVariant.ghost} onClick={onCancelClick} fullWidth>
          Cancel
        </Button>
      </div>
    </Modal>
  )
}

export { SignTransactionModal }
