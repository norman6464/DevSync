import { useEffect, type ReactNode } from 'react';
import { modalOverlayClass, modalContentClass } from '../../constants/styles';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  maxWidth?: string;
  children: ReactNode;
}

export default function Modal({
  isOpen,
  onClose,
  title,
  maxWidth = 'max-w-2xl',
  children,
}: ModalProps) {
  useEffect(() => {
    if (!isOpen) return;
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };
    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div
      data-testid="modal-overlay"
      className={modalOverlayClass}
      onClick={onClose}
    >
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby={title ? 'modal-title' : undefined}
        className={`${modalContentClass} ${maxWidth}`}
        onClick={(e) => e.stopPropagation()}
      >
        {title && (
          <h2 id="modal-title" className="text-xl font-semibold text-white mb-4">{title}</h2>
        )}
        {children}
      </div>
    </div>
  );
}
