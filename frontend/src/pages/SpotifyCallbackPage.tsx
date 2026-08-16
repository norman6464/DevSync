import { useEffect, useState, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { spotifyCallback } from '../api/spotify';
import { useAuthStore } from '../store/authStore';
import { PageLoader } from '../components/common';
import toast from 'react-hot-toast';

export default function SpotifyCallbackPage() {
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const loadUser = useAuthStore((s) => s.loadUser);
  const [asyncError, setError] = useState('');
  const hasProcessed = useRef(false);

  // パラメータ不足エラーは URL から一意に決まる派生値（effect で state に写さない）
  const code = searchParams.get('code');
  const state = searchParams.get('state');
  const error = asyncError || (!code || !state ? t('spotifyCallback.missingParams') : '');

  useEffect(() => {
    if (hasProcessed.current) return;
    hasProcessed.current = true;

    if (!code || !state) return;

    spotifyCallback(code, state)
      .then(async () => {
        await loadUser();
        toast.success(t('spotifyCallback.connectSuccess'));
        navigate('/settings');
      })
      .catch(() => {
        setError(t('spotifyCallback.connectFailed'));
      });
  }, [code, state, navigate, loadUser, t]);

  if (error) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-400 mb-4" role="alert">{error}</p>
          <button
            onClick={() => navigate('/settings')}
            className="px-4 py-2 bg-blue-600 text-white rounded-md"
          >
            {t('spotifyCallback.backToSettings')}
          </button>
        </div>
      </div>
    );
  }

  return <PageLoader fullHeight message={t('spotifyCallback.connecting')} />;
}
