import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { gitHubCallback } from '../api/github';
import { useAuthStore } from '../store/authStore';
import LoadingSpinner from '../components/common/LoadingSpinner';
import toast from 'react-hot-toast';

function parseStatePurpose(state: string): string {
  try {
    const payload = state.split('.')[1];
    const decoded = JSON.parse(atob(payload));
    return decoded.purpose || '';
  } catch {
    return '';
  }
}

export default function GitHubCallbackPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const loadUser = useAuthStore((s) => s.loadUser);
  const handleGitHubCallback = useAuthStore((s) => s.handleGitHubCallback);
  const [error, setError] = useState('');
  const [mode, setMode] = useState<'connect' | 'login' | ''>('');

  useEffect(() => {
    const code = searchParams.get('code');
    const state = searchParams.get('state');

    if (!code || !state) {
      setError('Missing OAuth parameters');
      return;
    }

    const purpose = parseStatePurpose(state);

    if (purpose === 'github_login') {
      setMode('login');
      handleGitHubCallback(code, state)
        .then(() => {
          toast.success('Logged in with GitHub!');
          navigate('/');
        })
        .catch(() => {
          setError('GitHub login failed');
        });
    } else {
      setMode('connect');
      gitHubCallback(code, state)
        .then(async () => {
          await loadUser();
          toast.success('GitHub connected successfully!');
          navigate('/settings');
        })
        .catch(() => {
          setError('Failed to connect GitHub');
        });
    }
  }, [searchParams, navigate, loadUser, handleGitHubCallback]);

  if (error) {
    const backPath = mode === 'login' ? '/login' : '/settings';
    const backLabel = mode === 'login' ? 'Back to Login' : 'Back to Settings';
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-400 mb-4">{error}</p>
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
    <div className="min-h-screen bg-gray-950 flex items-center justify-center">
      <div className="text-center">
        <LoadingSpinner />
        <p className="text-gray-400 mt-4">
          {mode === 'login' ? 'Logging in with GitHub...' : 'Connecting GitHub...'}
        </p>
      </div>
    </div>
  );
}
