interface Stat {
  label: string;
  value: number;
}

interface ProfileCardProps {
  name: string;
  title?: string;
  avatarUrl?: string;
  bio?: string;
  stats?: Stat[];
  onFollow?: () => void;
  isFollowing?: boolean;
  className?: string;
}

export default function ProfileCard({
  name,
  title,
  avatarUrl,
  bio,
  stats,
  onFollow,
  isFollowing = false,
  className = '',
}: ProfileCardProps) {
  return (
    <div className={`p-6 bg-gray-800/50 border border-gray-700 rounded-xl text-center ${className}`.trim()}>
      <div className="w-16 h-16 mx-auto rounded-full bg-gray-700 overflow-hidden flex items-center justify-center">
        {avatarUrl ? (
          <img src={avatarUrl} alt={name} className="w-full h-full object-cover" />
        ) : (
          <span className="text-xl font-bold text-gray-200">{name.charAt(0)}</span>
        )}
      </div>
      <h3 className="mt-3 text-lg font-semibold text-gray-200">{name}</h3>
      {title && <p className="text-sm text-gray-400">{title}</p>}
      {bio && <p className="mt-2 text-xs text-gray-500 leading-relaxed">{bio}</p>}
      {stats && stats.length > 0 && (
        <div className="flex justify-center gap-6 mt-4">
          {stats.map((stat) => (
            <div key={stat.label}>
              <p className="text-lg font-bold text-gray-200">{stat.value}</p>
              <p className="text-xs text-gray-500">{stat.label}</p>
            </div>
          ))}
        </div>
      )}
      {onFollow && (
        <button
          type="button"
          onClick={onFollow}
          className={`mt-4 px-6 py-2 text-sm font-medium rounded-lg transition-colors ${
            isFollowing
              ? 'bg-gray-700 text-gray-300 hover:bg-gray-600'
              : 'bg-blue-600 text-white hover:bg-blue-500'
          }`}
        >
          {isFollowing ? 'フォロー中' : 'フォロー'}
        </button>
      )}
    </div>
  );
}
