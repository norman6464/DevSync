import { ReactNode, useEffect } from 'react';
import { X } from 'lucide-react';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  children: ReactNode;
}

interface ModalHeaderProps {
  children: ReactNode;
}

interface ModalBodyProps {
  children: ReactNode;
}

interface ModalFooterProps {
  children: ReactNode;
}

interface ModalTitleProps {
  children: ReactNode;
}

function Modal({ isOpen, onClose, children }: ModalProps) {
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };

    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* オーバーレイ */}
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />

      {/* モーダルコンテンツ */}
      <div className="relative bg-gray-900 border border-gray-800 rounded-lg shadow-xl max-w-lg w-full max-h-[90vh] overflow-y-auto">
        {/* 閉じるボタン */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 p-1 text-gray-400 hover:text-white transition-colors rounded"
          aria-label="閉じる"
        >
          <X className="w-5 h-5" />
        </button>

        {children}
      </div>
    </div>
  );
}

function ModalHeader({ children }: ModalHeaderProps) {
  return <div className="px-6 pt-6 pb-4">{children}</div>;
}

function ModalBody({ children }: ModalBodyProps) {
  return <div className="px-6 py-4">{children}</div>;
}

function ModalFooter({ children }: ModalFooterProps) {
  return (
    <div className="px-6 py-4 border-t border-gray-800 flex items-center justify-end gap-3">
      {children}
    </div>
  );
}

function ModalTitle({ children }: ModalTitleProps) {
  return <h2 className="text-xl font-semibold text-white">{children}</h2>;
}

Modal.Header = ModalHeader;
Modal.Body = ModalBody;
Modal.Footer = ModalFooter;
Modal.Title = ModalTitle;

export default Modal;
