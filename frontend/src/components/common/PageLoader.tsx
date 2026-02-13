import { useTranslation } from 'react-i18next';
import LoadingSpinner from './LoadingSpinner';

interface PageLoaderProps {
  message?: string;
  showMessage?: boolean;
  fullHeight?: boolean;
}

export default function PageLoader({
  message,
  showMessage = true,
  fullHeight = false,
}: PageLoaderProps) {
  const { t } = useTranslation();
  const displayMessage = message || t('common.loading');

  return (
    <div
      className={`flex flex-col items-center justify-center ${
        fullHeight ? 'min-h-screen' : 'py-12'
      }`}
    >
      <div className="w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full animate-spin" />
      {showMessage && (
        <p className="mt-4 text-gray-400 text-sm">{displayMessage}</p>
      )}
    </div>
  );
}
