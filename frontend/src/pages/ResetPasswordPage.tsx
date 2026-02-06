import { useState } from 'react';
import { Link, useSearchParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { resetPassword } from '../api/auth';

export default function ResetPasswordPage() {
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const token = searchParams.get('token') || '';

  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (password !== confirmPassword) {
      setError(t('accountManagement.passwordsMismatch'));
      return;
    }

    if (password.length < 6) {
      setError(t('accountManagement.passwordTooShort'));
      return;
    }

    setLoading(true);

    try {
      await resetPassword(token, password);
      setSuccess(true);
      setTimeout(() => navigate('/login'), 3000);
    } catch (err) {
      setError(t('accountManagement.resetFailed'));
    } finally {
      setLoading(false);
    }
  };

  if (!token) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center px-4">
        <div className="max-w-md w-full">
          <div className="bg-gray-900 border border-gray-800 rounded-lg p-6 text-center">
            <p className="text-red-400 mb-4">{t('accountManagement.invalidToken')}</p>
            <Link
              to="/forgot-password"
              className="text-blue-400 hover:text-blue-300"
            >
              {t('accountManagement.requestNewReset')}
            </Link>
          </div>
        </div>
      </div>
    );
  }

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
            {t('accountManagement.resetPassword')}
          </h1>
          <p className="text-gray-400 text-sm mb-6">
            {t('accountManagement.resetPasswordDesc')}
          </p>

          {success ? (
            <div className="space-y-4">
              <div className="bg-green-500/10 border border-green-500/30 text-green-400 px-4 py-3 rounded-lg">
                {t('accountManagement.passwordResetSuccess')}
              </div>
              <p className="text-gray-400 text-sm text-center">
                {t('accountManagement.redirectingToLogin')}
              </p>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              {error && (
                <div className="bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-3 rounded-lg text-sm">
                  {error}
                </div>
              )}

              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1.5">
                  {t('accountManagement.newPassword')}
                </label>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500"
                  placeholder="••••••••"
                  required
                  minLength={6}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1.5">
                  {t('auth.confirmPassword')}
                </label>
                <input
                  type="password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500"
                  placeholder="••••••••"
                  required
                  minLength={6}
                />
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-600/50 text-white rounded-lg font-medium transition-colors"
              >
                {loading ? t('common.loading') : t('accountManagement.resetPassword')}
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
