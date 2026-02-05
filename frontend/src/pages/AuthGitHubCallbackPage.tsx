import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import LoadingSpinner from '../components/common/LoadingSpinner';
import toast from 'react-hot-toast';

export default function AuthGitHubCallbackPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const handleGitHubCallback = useAuthStore((s) => s.handleGitHubCallback);
  const [error, setError] = useState('');

  useEffect(() => {
    const code = searchParams.get('code');
    const state = searchParams.get('state');

    if (!code || !state) {
      setError('Missing OAuth parameters');
      return;
    }

    handleGitHubCallback(code, state)
      .then(() => {
        toast.success('Logged in with GitHub!');
        navigate('/');
      })
      .catch(() => {
        setError('GitHub login failed');
      });
  }, [searchParams, navigate, handleGitHubCallback]);

  if (error) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-400 mb-4">{error}</p>
          <button
            onClick={() => navigate('/login')}
            className="px-4 py-2 bg-blue-600 text-white rounded-md"
          >
            Back to Login
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-950 flex items-center justify-center">
      <div className="text-center">
        <LoadingSpinner />
        <p className="text-gray-400 mt-4">Logging in with GitHub...</p>
      </div>
    </div>
  );
}
