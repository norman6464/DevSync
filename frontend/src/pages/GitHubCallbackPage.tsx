import { useEffect, useState, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { gitHubCallback } from '../api/github';
import { useAuthStore } from '../store/authStore';
import { PageLoader } from '../components/common';
import toast from 'react-hot-toast';

function parseStatePurpose(state: string): string {
  try {
    const parts = state.split('.');
    if (parts.length < 2 || !parts[1]) return '';
    const decoded = JSON.parse(atob(parts[1]));
    return typeof decoded.purpose === 'string' ? decoded.purpose : '';
  } catch {
    return '';
  }
}

export default function GitHubCallbackPage() {
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const loadUser = useAuthStore((s) => s.loadUser);
  const handleGitHubCallback = useAuthStore((s) => s.handleGitHubCallback);
  const [asyncError, setAsyncError] = useState('');
  const hasProcessed = useRef(false);

  // mode とパラメータ不足エラーは URL から一意に決まる派生値（effect で state に写さない）
  const code = searchParams.get('code');
  const state = searchParams.get('state');
  const mode: 'connect' | 'login' = state && parseStatePurpose(state) === 'github_login' ? 'login' : 'connect';
  const error = asyncError || (!code || !state ? t('githubCallback.missingOAuthParams') : '');

  useEffect(() => {
    // 二重実行を防ぐ
    if (hasProcessed.current) return;
    hasProcessed.current = true;

    if (!code || !state) return;

    if (parseStatePurpose(state) === 'github_login') {
      handleGitHubCallback(code, state)
        .then(() => {
          toast.success(t('githubCallback.loginSuccess'));
          navigate('/');
        })
        .catch(() => {
          setAsyncError(t('githubCallback.loginFailed'));
        });
    } else {
      gitHubCallback(code, state)
        .then(async () => {
          await loadUser();
          toast.success(t('githubCallback.connectSuccess'));
          const onboardingRedirect = sessionStorage.getItem('onboarding_redirect');
          if (onboardingRedirect === 'true') {
            sessionStorage.removeItem('onboarding_redirect');
            navigate('/onboarding');
          } else {
            navigate('/settings');
          }
        })
        .catch(() => {
          setAsyncError(t('githubCallback.connectFailed'));
        });
    }
  }, [code, state, navigate, loadUser, handleGitHubCallback, t]);

  if (error) {
    const backPath = mode === 'login' ? '/login' : '/settings';
    const backLabel = mode === 'login' ? t('githubCallback.backToLogin') : t('githubCallback.backToSettings');
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-400 mb-4" role="alert">{error}</p>
          <button
            onClick={() => navigate(backPath)}
            className="px-4 py-2 bg-blue-600 text-white rounded-md"
          >
            {backLabel}
          </button>
        </div>
      </div>
    );
  }

  return (
    <PageLoader
      fullHeight
      message={mode === 'login' ? t('githubCallback.loggingIn') : t('githubCallback.connecting')}
    />
  );
}
