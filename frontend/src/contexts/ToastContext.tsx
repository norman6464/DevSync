import { createContext, useContext, useState, useCallback, type ReactNode } from 'react';

interface Toast {
  id: number;
  message: string;
  icon?: string;
  color?: string;
}

interface ToastContextValue {
  showToast: (message: string, options?: { icon?: string; color?: string }) => void;
}

const ToastContext = createContext<ToastContextValue | null>(null);

let nextId = 0;

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const showToast = useCallback(
    (message: string, options?: { icon?: string; color?: string }) => {
      const id = nextId++;
      setToasts((prev) => [...prev.slice(-2), { id, message, ...options }]);

      setTimeout(() => {
        setToasts((prev) => prev.filter((t) => t.id !== id));
      }, 5000);
    },
    []
  );

  const removeToast = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  return (
    <ToastContext.Provider value={{ showToast }}>
      {children}
      {/* Toast Container */}
      {toasts.length > 0 && (
        <div className="fixed top-20 right-4 z-[100] flex flex-col gap-2 pointer-events-none">
          {toasts.map((toast) => (
            <div
              key={toast.id}
              className="pointer-events-auto animate-slide-in-right flex items-center gap-3 px-4 py-3 bg-gray-800 border border-gray-700 rounded-lg shadow-xl max-w-sm"
            >
              {toast.icon && <span className="text-xl shrink-0">{toast.icon}</span>}
              <p className={`text-sm font-medium ${toast.color || 'text-white'}`}>
                {toast.message}
              </p>
              <button
                onClick={() => removeToast(toast.id)}
                className="ml-auto shrink-0 text-gray-400 hover:text-white transition-colors"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          ))}
        </div>
      )}
    </ToastContext.Provider>
  );
}

export function useToast() {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error('useToast must be used within ToastProvider');
  return ctx;
}
