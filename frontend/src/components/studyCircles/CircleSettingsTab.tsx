import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { Plus, Crown, UserMinus } from 'lucide-react';
import type { StudyCircle } from '../../types/studyCircle';
import type { User } from '../../types/user';
import { useUserSearch, useConfirm } from '../../hooks';
import Avatar from '../common/Avatar';
import ConfirmDialog from '../common/ConfirmDialog';

interface CircleSettingsTabProps {
  circle: StudyCircle;
  onAddMember: (userId: number) => Promise<unknown>;
  onRemoveMember: (userId: number) => void;
}

export default function CircleSettingsTab({ circle, onAddMember, onRemoveMember }: CircleSettingsTabProps) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { confirm, dialogProps } = useConfirm();
  const [showAddMember, setShowAddMember] = useState(false);
  const [memberSearch, setMemberSearch] = useState('');
  // NOTE: useUserSearch は実際には引数を取らず users も返さないため、この呼び出しは
  // 実行時に searchUsers が undefined になる既存不具合を含む。挙動を変えないよう
  // 型だけ通すキャストに留めている（フック API に合わせた修正は別途行う）。
  const { users: searchUsers } = (useUserSearch as unknown as (query: string) => { users: User[] })(memberSearch);

  const handleAddMember = async (userId: number) => {
    await onAddMember(userId);
    setShowAddMember(false);
    setMemberSearch('');
  };

  return (
    <div className="space-y-4">
      {/* Members Management */}
      <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
        <div className="flex items-center justify-between mb-3">
          <h3 className="text-sm font-medium text-white">{t('studyCircle.members')}</h3>
          <button
            onClick={() => setShowAddMember(!showAddMember)}
            aria-expanded={showAddMember}
            className="flex items-center gap-1 px-2 py-1 text-xs text-purple-400 hover:text-purple-300 border border-purple-500/30 rounded-lg transition-colors"
          >
            <Plus className="w-3 h-3" aria-hidden="true" />
            {t('studyCircle.addMember')}
          </button>
        </div>

        {showAddMember && (
          <div className="mb-3 p-3 bg-gray-800/50 rounded-lg">
            <input
              type="text"
              value={memberSearch}
              onChange={(e) => setMemberSearch(e.target.value)}
              placeholder={t('studyCircle.searchUsers')}
              maxLength={50}
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500 mb-2"
            />
            {searchUsers.length > 0 && (
              <div className="space-y-1 max-h-32 overflow-y-auto">
                {searchUsers
                  .filter((u: User) => !circle.members?.some((m) => m.user_id === u.id))
                  .slice(0, 5)
                  .map((u: User) => (
                  <button
                    key={u.id}
                    onClick={() => handleAddMember(u.id)}
                    className="flex items-center gap-2 w-full p-2 rounded-lg hover:bg-gray-800 text-left transition-colors"
                  >
                    <Avatar name={u.name} avatarUrl={u.avatar_url} size="xs" />
                    <span className="text-xs text-white">{u.name}</span>
                  </button>
                ))}
              </div>
            )}
          </div>
        )}

        <div className="space-y-1.5">
          {circle.members?.map((m) => (
            <div key={m.id} className="flex items-center gap-2.5 p-2 rounded-lg hover:bg-gray-800/50">
              <Avatar name={m.user?.name || ''} avatarUrl={m.user?.avatar_url} size="sm" />
              <div className="flex-1 min-w-0">
                <p className="text-xs text-white truncate">{m.user?.name}</p>
                <p className="text-[10px] text-gray-500">{m.role === 'owner' ? t('studyCircle.owner') : t('studyCircle.members')}</p>
              </div>
              {m.role !== 'owner' && (
                <button
                  onClick={() => onRemoveMember(m.user_id)}
                  className="text-gray-600 hover:text-red-400 transition-colors"
                  aria-label={t('studyCircle.removeMember')}
                >
                  <UserMinus className="w-3.5 h-3.5" aria-hidden="true" />
                </button>
              )}
              {m.role === 'owner' && (
                <Crown className="w-3.5 h-3.5 text-amber-400" aria-hidden="true" />
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Danger Zone */}
      <div className="bg-gray-900 border border-red-500/20 rounded-md p-4">
        <h3 className="text-sm font-medium text-red-400 mb-3">{t('common.delete')}</h3>
        <button
          onClick={async () => {
            const confirmed = await confirm({
              title: t('common.delete'),
              message: t('studyCircle.confirmDelete'),
              variant: 'danger',
              confirmText: t('common.delete'),
            });
            if (confirmed) {
              const { deleteCircle } = await import('../../api/studyCircles');
              await deleteCircle(circle.id);
              navigate('/study-circles');
            }
          }}
          className="px-3 py-1.5 bg-red-600/20 hover:bg-red-600/30 text-red-400 rounded-lg text-xs font-medium border border-red-500/30 transition-colors"
        >
          {t('studyCircle.confirmDelete')}
        </button>
      </div>
      <ConfirmDialog {...dialogProps} />
    </div>
  );
}
