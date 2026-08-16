import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { requestPasswordReset } from '../api/auth';
import { inputClass, labelClass } from '../constants/styles';

export default function ForgotPasswordPage() {
  const { t } = useTranslation();
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      await requestPasswordReset(email);
      setSuccess(true);
    } catch {
      setError(t('accountManagement.resetRequestFailed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-950 flex items-center justify-center px-4">
      <div className="max-w-md w-full">
        <div className="text-center mb-8">
          <Link to="/" className="inline-flex items-center gap-2">
            <svg className="w-10 h-10 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M16 18l2-2-2-2" />
              <path d="M8 6L6 8l2 2" />
              <path d="M14.5 4l-5 16" />
            </svg>
            <span className="text-2xl font-bold text-white">DevSync</span>
          </Link>
        </div>

        <div className="bg-gray-900 border border-gray-800 rounded-lg p-6">
          <h1 className="text-xl font-semibold text-white mb-2">
            {t('accountManagement.forgotPassword')}
          </h1>
          <p className="text-gray-400 text-sm mb-6">
            {t('accountManagement.forgotPasswordDesc')}
          </p>

          {success ? (
            <div className="space-y-4">
              <div className="bg-green-500/10 border border-green-500/30 text-green-400 px-4 py-3 rounded-lg">
                {t('accountManagement.resetEmailSent')}
              </div>

              <Link
                to="/login"
                className="block w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium text-center transition-colors"
              >
                {t('auth.login')}
              </Link>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              {error && (
                <div className="bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-3 rounded-lg text-sm">
                  {error}
                </div>
              )}

              <div>
                <label className={labelClass}>
                  {t('auth.email')}
                </label>
                <input
                  type="email"
                  autoComplete="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className={inputClass}
                  placeholder={t('auth.emailPlaceholder')}
                  required
                  maxLength={255}
                />
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-600/50 text-white rounded-lg font-medium transition-colors"
              >
                {loading ? t('common.loading') : t('accountManagement.sendResetLink')}
              </button>
            </form>
          )}

          <div className="mt-6 text-center">
            <Link to="/login" className="text-blue-400 hover:text-blue-300 text-sm">
              {t('accountManagement.backToLogin')}
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
