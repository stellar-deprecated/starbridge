import ReactModal from 'react-modal'

import classNames from 'classnames'

import styles from './styles.module.scss'

export interface IModalProps {
  className?: string
  children: React.ReactNode
  isOpen: boolean
  setModalOpen: (isOpen: boolean) => void
}

const defaultStyles = {
  overlay: {
    backgroundColor: 'var(--color-blackish-with-opacity)',
  },
}

const Modal = ({
  className,
  children,
  isOpen = false,
  setModalOpen,
}: IModalProps): JSX.Element => {
  return (
    <ReactModal
      isOpen={isOpen}
      className={classNames(styles.modal, className)}
      style={defaultStyles}
      closeTimeoutMS={300}
      shouldCloseOnOverlayClick
      shouldCloseOnEsc
      onRequestClose={(): void => setModalOpen(false)}
    >
      {children}
    </ReactModal>
  )
}

export { Modal }
