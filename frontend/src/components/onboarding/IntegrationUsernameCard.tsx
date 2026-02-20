import { useTranslation } from 'react-i18next';
import { CheckCircle } from 'lucide-react';
import { inputClass } from '../../constants/styles';

interface IntegrationUsernameCardProps {
  icon: React.ReactNode;
  serviceName: string;
  description: string;
  connectedUsername?: string;
  username: string;
  onUsernameChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  placeholder: string;
  connecting: boolean;
  onConnect: () => void;
  buttonClassName: string;
}

export default function IntegrationUsernameCard({
  icon,
  serviceName,
  description,
  connectedUsername,
  username,
  onUsernameChange,
  placeholder,
  connecting,
  onConnect,
  buttonClassName,
}: IntegrationUsernameCardProps) {
  const { t } = useTranslation();

  return (
    <div className="p-4 bg-gray-800/50 rounded-lg border border-gray-700">
      <div className="flex items-center gap-3 mb-3">
        {icon}
        <div>
          <h3 className="text-sm font-medium text-white">{serviceName}</h3>
          <p className="text-xs text-gray-400">{description}</p>
        </div>
      </div>
      {connectedUsername ? (
        <div className="flex items-center gap-2 text-green-400 text-sm">
          <CheckCircle className="w-4 h-4" />
          <span>{t('settings.connected')} - @{connectedUsername}</span>
        </div>
      ) : (
        <div className="flex gap-2">
          <input
            type="text"
            value={username}
            onChange={onUsernameChange}
            placeholder={placeholder}
            maxLength={50}
            className={`${inputClass} flex-1`}
          />
          <button
            onClick={onConnect}
            disabled={connecting || !username.trim()}
            className={`${buttonClassName} disabled:opacity-50 whitespace-nowrap`}
          >
            {connecting ? t('common.loading') : t('settings.connect')}
          </button>
        </div>
      )}
    </div>
  );
}
