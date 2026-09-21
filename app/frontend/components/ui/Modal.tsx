import { ReactNode } from 'react';

export type ModalProps = {
  children: ReactNode;
  /** パネルの幅。既定は w-96 */
  className?: string;
};

// オーバーレイ + 中央寄せパネルの共通の外枠。
const Modal = ({ children, className = 'w-96' }: ModalProps) => (
  <div className='fixed inset-0 z-50 flex items-center justify-center bg-black/50'>
    <div className={`bg-surface rounded-lg shadow-lg mx-4 ${className}`.trim()}>
      {children}
    </div>
  </div>
);

export default Modal;
