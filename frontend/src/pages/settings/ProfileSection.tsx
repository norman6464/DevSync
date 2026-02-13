import { useTranslation } from 'react-i18next';
import { inputClass } from '../../constants/styles';

interface Props {
  name: string;
  setName: (v: string) => void;
  bio: string;
  setBio: (v: string) => void;
  avatarUrl: string;
  setAvatarUrl: (v: string) => void;
  saving: boolean;
  onSubmit: (e: React.FormEvent) => void;
}

export default function ProfileSection({ name, setName, bio, setBio, avatarUrl, setAvatarUrl, saving, onSubmit }: Props) {
  const { t } = useTranslation();

  return (
    <form onSubmit={onSubmit} className="bg-gray-900 border border-gray-800 rounded-md overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-800">
        <h2 className="text-base font-semibold">{t('settings.profile')}</h2>
      </div>
      <div className="p-6 space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1.5">{t('settings.name')}</label>
          <input type="text" value={name} onChange={(e) => setName(e.target.value)} className={inputClass} />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1.5">{t('settings.bio')}</label>
          <textarea value={bio} onChange={(e) => setBio(e.target.value)} rows={3} placeholder="Tell us about yourself" className={`${inputClass} resize-none`} />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1.5">{t('settings.avatar')}</label>
          <input type="text" value={avatarUrl} onChange={(e) => setAvatarUrl(e.target.value)} placeholder="https://..." className={inputClass} />
        </div>
      </div>
      <div className="px-6 py-4 border-t border-gray-800 flex justify-end">
        <button
          type="submit"
          disabled={saving}
          className="px-5 py-2 bg-gray-700 hover:bg-gray-600 disabled:opacity-50 text-white rounded-lg font-medium text-sm transition-colors"
        >
          {saving ? t('common.loading') : t('common.save')}
        </button>
      </div>
    </form>
  );
}
