import { useTranslation } from 'react-i18next';
import { inputClass, buttonSecondaryClass, sectionContainerClass, labelClass, textareaClass } from '../../constants/styles';

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
    <form onSubmit={onSubmit} className={sectionContainerClass}>
      <div className="px-6 py-4 border-b border-gray-800">
        <h2 className="text-base font-semibold">{t('settings.profile')}</h2>
      </div>
      <div className="p-6 space-y-4">
        <div>
          <label htmlFor="profile-name" className={labelClass}>{t('settings.name')}</label>
          <input id="profile-name" type="text" value={name} onChange={(e) => setName(e.target.value)} className={inputClass} />
        </div>
        <div>
          <label htmlFor="profile-bio" className={labelClass}>{t('settings.bio')}</label>
          <textarea id="profile-bio" value={bio} onChange={(e) => setBio(e.target.value)} rows={3} placeholder={t('settings.bioPlaceholder')} className={textareaClass} />
        </div>
        <div>
          <label htmlFor="profile-avatar" className={labelClass}>{t('settings.avatar')}</label>
          <input id="profile-avatar" type="text" value={avatarUrl} onChange={(e) => setAvatarUrl(e.target.value)} placeholder="https://..." className={inputClass} />
        </div>
      </div>
      <div className="px-6 py-4 border-t border-gray-800 flex justify-end">
        <button
          type="submit"
          disabled={saving}
          className={`${buttonSecondaryClass} px-5 font-medium text-sm disabled:opacity-50`}
        >
          {saving ? t('common.loading') : t('common.save')}
        </button>
      </div>
    </form>
  );
}
