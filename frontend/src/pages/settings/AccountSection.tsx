import { useTranslation } from 'react-i18next';
import { inputClass } from '../../constants/styles';
import type { User } from '../../types/user';

interface Props {
  user: User;
  // Email
  emailWeeklyReport: boolean;
  setEmailWeeklyReport: (v: boolean) => void;
  emailLanguage: string;
  setEmailLanguage: (v: string) => void;
  savingEmail: boolean;
  onSaveEmailPreferences: () => void;
  // Delete
  showDeleteModal: boolean;
  setShowDeleteModal: (v: boolean) => void;
  deletePassword: string;
  setDeletePassword: (v: string) => void;
  deleting: boolean;
  onDeleteAccount: () => void;
}

export default function AccountSection(props: Props) {
  const { t } = useTranslation();

  return (
    <>
      {/* Email Notifications */}
      <div className="bg-gray-900 border border-gray-800 rounded-md overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-800">
          <h2 className="text-base font-semibold">{t('settings.emailNotifications')}</h2>
        </div>
        <div className="p-6 space-y-5">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-200">{t('settings.emailWeeklyReport')}</p>
              <p className="text-xs text-gray-500 mt-0.5">{t('settings.emailWeeklyReportDesc')}</p>
            </div>
            <button
              type="button"
              onClick={() => props.setEmailWeeklyReport(!props.emailWeeklyReport)}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                props.emailWeeklyReport ? 'bg-blue-600' : 'bg-gray-700'
              }`}
            >
              <span
                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                  props.emailWeeklyReport ? 'translate-x-6' : 'translate-x-1'
                }`}
              />
            </button>
          </div>
          <div className="text-xs">
            <span className={props.emailWeeklyReport ? 'text-blue-400' : 'text-gray-500'}>
              {props.emailWeeklyReport ? t('settings.emailEnabled') : t('settings.emailDisabled')}
            </span>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1.5">{t('settings.emailLanguage')}</label>
            <p className="text-xs text-gray-500 mb-2">{t('settings.emailLanguageDesc')}</p>
            <select
              value={props.emailLanguage}
              onChange={(e) => props.setEmailLanguage(e.target.value)}
              className={inputClass}
            >
              <option value="ja">日本語</option>
              <option value="en">English</option>
              <option value="ko">한국어</option>
              <option value="zh-CN">简体中文</option>
              <option value="zh-TW">繁體中文</option>
              <option value="es">Español</option>
              <option value="fr">Français</option>
              <option value="de">Deutsch</option>
              <option value="pt">Português</option>
              <option value="ru">Русский</option>
            </select>
          </div>
        </div>
        <div className="px-6 py-4 border-t border-gray-800 flex justify-end">
          <button
            type="button"
            onClick={props.onSaveEmailPreferences}
            disabled={props.savingEmail}
            className="px-5 py-2 bg-gray-700 hover:bg-gray-600 disabled:opacity-50 text-white rounded-lg font-medium text-sm transition-colors"
          >
            {props.savingEmail ? t('common.loading') : t('common.save')}
          </button>
        </div>
      </div>

      {/* Danger Zone */}
      <div className="bg-gray-900 border border-red-500/30 rounded-md overflow-hidden">
        <div className="px-6 py-4 border-b border-red-500/30 bg-red-500/5">
          <h2 className="text-base font-semibold text-red-400">{t('accountManagement.dangerZone')}</h2>
        </div>
        <div className="p-6">
          <p className="text-gray-400 text-sm mb-4">
            {t('accountManagement.deleteWarning')}
          </p>
          <button
            onClick={() => props.setShowDeleteModal(true)}
            className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg text-sm font-medium transition-colors"
          >
            {t('accountManagement.deleteAccount')}
          </button>
        </div>
      </div>

      {/* Delete Confirmation Modal */}
      {props.showDeleteModal && (
        <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50 px-4">
          <div className="bg-gray-900 border border-gray-700 rounded-md max-w-md w-full p-6">
            <h3 className="text-lg font-semibold text-white mb-2">
              {t('accountManagement.confirmDelete')}
            </h3>
            <p className="text-gray-400 text-sm mb-4">
              {t('accountManagement.deleteConfirmText')}
            </p>

            {!props.user.github_connected && (
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-300 mb-1.5">
                  {t('auth.password')}
                </label>
                <input
                  type="password"
                  value={props.deletePassword}
                  onChange={(e) => props.setDeletePassword(e.target.value)}
                  className={inputClass}
                  placeholder="••••••••"
                />
              </div>
            )}

            <div className="flex gap-3 justify-end">
              <button
                onClick={() => {
                  props.setShowDeleteModal(false);
                  props.setDeletePassword('');
                }}
                className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg text-sm font-medium transition-colors"
              >
                {t('common.cancel')}
              </button>
              <button
                onClick={props.onDeleteAccount}
                disabled={props.deleting}
                className="px-4 py-2 bg-red-600 hover:bg-red-700 disabled:opacity-50 text-white rounded-lg text-sm font-medium transition-colors"
              >
                {props.deleting ? t('common.loading') : t('accountManagement.deleteAccount')}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
