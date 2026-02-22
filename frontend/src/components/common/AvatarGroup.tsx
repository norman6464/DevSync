interface AvatarUser {
  name: string;
  image?: string;
}

const sizeMap = {
  sm: 'w-8 h-8 text-xs',
  md: 'w-10 h-10 text-sm',
  lg: 'w-12 h-12 text-base',
};

interface AvatarGroupProps {
  users: AvatarUser[];
  max?: number;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export default function AvatarGroup({
  users,
  max,
  size = 'md',
  className = '',
}: AvatarGroupProps) {
  const visibleUsers = max ? users.slice(0, max) : users;
  const remaining = max ? users.length - max : 0;

  return (
    <div className={`flex items-center ${className}`.trim()}>
      {visibleUsers.map((user, index) => (
        <div
          key={index}
          className={`${sizeMap[size]} rounded-full border-2 border-gray-900 flex items-center justify-center overflow-hidden ${index > 0 ? '-ml-2' : ''}`}
        >
          {user.image ? (
            <img
              src={user.image}
              alt={user.name}
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full bg-gray-700 flex items-center justify-center text-gray-300 font-medium">
              {user.name.charAt(0).toUpperCase()}
            </div>
          )}
        </div>
      ))}
      {remaining > 0 && (
        <div
          className={`${sizeMap[size]} rounded-full border-2 border-gray-900 -ml-2 bg-gray-700 flex items-center justify-center text-gray-300 font-medium`}
        >
          +{remaining}
        </div>
      )}
    </div>
  );
}
