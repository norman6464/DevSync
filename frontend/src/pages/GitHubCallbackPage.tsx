import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { gitHubCallback } from '../api/github';
import { useAuthStore } from '../store/authStore';
import LoadingSpinner from '../components/common/LoadingSpinner';
import toast from 'react-hot-toast';

export default function GitHubCallbackPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const loadUser = useAuthStore((s) => s.loadUser);
  const [error, setError] = useState('');

  useEffect(() => {
    const code = searchParams.get('code');
    const state = searchParams.get('state');

    if (!code || !state) {
      setError('Missing OAuth parameters');
      return;
    }

    gitHubCallback(code, state)
      .then(async () => {
        await loadUser();
        toast.success('GitHub connected successfully!');
        navigate('/settings');
      })
      .catch(() => {
        setError('Failed to connect GitHub');
      });
  }, [searchParams, navigate, loadUser]);

  if (error) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-400 mb-4">{error}</p>
          <button
            onClick={() => navigate('/settings')}
            className="px-4 py-2 bg-blue-600 text-white rounded-md"
          >
            Back to Settings
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-950 flex items-center justify-center">
      <div className="text-center">
        <LoadingSpinner />
        <p className="text-gray-400 mt-4">Connecting GitHub...</p>
      </div>
    </div>
  );
}
