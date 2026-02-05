import { useState } from 'react';
import { useAuthStore } from '../store/authStore';
import { updateUser } from '../api/users';
import { getGitHubConnectURL, disconnectGitHub, syncGitHub } from '../api/github';
import toast from 'react-hot-toast';

export default function SettingsPage() {
  const { user, setUser } = useAuthStore();
  const [name, setName] = useState(user?.name || '');
  const [bio, setBio] = useState(user?.bio || '');
  const [avatarUrl, setAvatarUrl] = useState(user?.avatar_url || '');
  const [saving, setSaving] = useState(false);
  const [syncing, setSyncing] = useState(false);

  if (!user) return null;

  const handleSaveProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      const { data } = await updateUser(user.id, { name, bio, avatar_url: avatarUrl });
      setUser(data);
      toast.success('Profile updated');
    } catch {
      toast.error('Failed to update profile');
    } finally {
      setSaving(false);
    }
  };

  const handleConnectGitHub = async () => {
    try {
      const { data } = await getGitHubConnectURL();
      window.location.href = data.url;
    } catch {
      toast.error('Failed to initiate GitHub connection');
    }
  };

  const handleDisconnectGitHub = async () => {
    try {
      await disconnectGitHub();
      setUser({ ...user, github_connected: false, github_username: '' });
      toast.success('GitHub disconnected');
    } catch {
      toast.error('Failed to disconnect GitHub');
    }
  };

  const handleSyncGitHub = async () => {
    setSyncing(true);
    try {
      await syncGitHub();
      toast.success('GitHub data synced');
    } catch {
      toast.error('Failed to sync GitHub data');
    } finally {
      setSyncing(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto space-y-8">
      <h1 className="text-2xl font-bold">Settings</h1>

      <form onSubmit={handleSaveProfile} className="bg-gray-900 rounded-lg p-6 space-y-4">
        <h2 className="text-lg font-semibold">Profile</h2>
        <div>
          <label className="block text-sm text-gray-400 mb-1">Name</label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white focus:outline-none focus:border-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm text-gray-400 mb-1">Bio</label>
          <textarea
            value={bio}
            onChange={(e) => setBio(e.target.value)}
            rows={3}
            className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white focus:outline-none focus:border-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm text-gray-400 mb-1">Avatar URL</label>
          <input
            type="text"
            value={avatarUrl}
            onChange={(e) => setAvatarUrl(e.target.value)}
            className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white focus:outline-none focus:border-blue-500"
          />
        </div>
        <button
          type="submit"
          disabled={saving}
          className="px-6 py-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white rounded-md font-medium transition-colors"
        >
          {saving ? 'Saving...' : 'Save'}
        </button>
      </form>

      <div className="bg-gray-900 rounded-lg p-6 space-y-4">
        <h2 className="text-lg font-semibold">GitHub Integration</h2>
        {user.github_connected ? (
          <div className="space-y-3">
            <p className="text-green-400">
              Connected as @{user.github_username}
            </p>
            <div className="flex gap-2">
              <button
                onClick={handleSyncGitHub}
                disabled={syncing}
                className="px-4 py-2 bg-gray-700 hover:bg-gray-600 disabled:opacity-50 text-white rounded-md text-sm transition-colors"
              >
                {syncing ? 'Syncing...' : 'Sync Data'}
              </button>
              <button
                onClick={handleDisconnectGitHub}
                className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-md text-sm transition-colors"
              >
                Disconnect
              </button>
            </div>
          </div>
        ) : (
          <button
            onClick={handleConnectGitHub}
            className="px-6 py-2 bg-gray-800 hover:bg-gray-700 text-white rounded-md font-medium border border-gray-600 transition-colors"
          >
            Connect GitHub
          </button>
        )}
      </div>
    </div>
  );
}
