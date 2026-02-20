import { useTranslation } from 'react-i18next';
import { inputClass, buttonSecondaryClass, buttonDangerClass, sectionContainerClass, labelClass, selectClass } from '../../constants/styles';
import { Modal } from '../../components/common';
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
      <div className={sectionContainerClass}>
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
            <label htmlFor="account-email-language" className={labelClass}>{t('settings.emailLanguage')}</label>
            <p className="text-xs text-gray-500 mb-2">{t('settings.emailLanguageDesc')}</p>
            <select
              id="account-email-language"
              value={props.emailLanguage}
              onChange={(e) => props.setEmailLanguage(e.target.value)}
              className={`${selectClass} w-full`}
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
            className={`${buttonSecondaryClass} px-5 font-medium text-sm disabled:opacity-50`}
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
            className={`${buttonDangerClass} text-sm`}
          >
            {t('accountManagement.deleteAccount')}
          </button>
        </div>
      </div>

      {/* Delete Confirmation Modal */}
      <Modal
        isOpen={props.showDeleteModal}
        onClose={() => { props.setShowDeleteModal(false); props.setDeletePassword(''); }}
        title={t('accountManagement.confirmDelete')}
        maxWidth="max-w-md"
      >
        <p className="text-gray-400 text-sm mb-4">
          {t('accountManagement.deleteConfirmText')}
        </p>

        {!props.user.github_connected && (
          <div className="mb-4">
            <label htmlFor="account-delete-password" className={labelClass}>
              {t('auth.password')}
            </label>
            <input
              id="account-delete-password"
              type="password"
              autoComplete="current-password"
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
            className={`${buttonSecondaryClass} text-sm font-medium`}
          >
            {t('common.cancel')}
          </button>
          <button
            onClick={props.onDeleteAccount}
            disabled={props.deleting}
            className={`${buttonDangerClass} text-sm disabled:opacity-50`}
          >
            {props.deleting ? t('common.loading') : t('accountManagement.deleteAccount')}
          </button>
        </div>
      </Modal>
    </>
  );
}
