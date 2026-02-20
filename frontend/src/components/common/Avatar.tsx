import { sanitizeUrl } from '../../utils/url';

interface AvatarProps {
  name: string;
  avatarUrl?: string;
  size?: 'xs' | 'sm' | 'md' | 'lg';
  isOnline?: boolean;
}

const dotSizeClasses = {
  xs: 'w-1.5 h-1.5',
  sm: 'w-2 h-2',
  md: 'w-2.5 h-2.5',
  lg: 'w-3.5 h-3.5',
};

const sizeClasses = {
  xs: 'w-6 h-6 text-xs',
  sm: 'w-8 h-8 text-sm',
  md: 'w-10 h-10 text-base',
  lg: 'w-16 h-16 text-xl',
};

export default function Avatar({ name, avatarUrl, size = 'md', isOnline }: AvatarProps) {
  const safeUrl = sanitizeUrl(avatarUrl);

  const indicator = isOnline != null && (
    <span
      className={`absolute bottom-0 right-0 ${dotSizeClasses[size]} rounded-full ring-2 ring-gray-900 ${isOnline ? 'bg-green-500' : 'bg-gray-500'}`}
      aria-hidden="true"
    />
  );

  if (safeUrl) {
    return (
      <span className="relative inline-block">
        <img
          src={safeUrl}
          alt={name}
          referrerPolicy="no-referrer"
          className={`${sizeClasses[size]} rounded-full object-cover`}
        />
        {indicator}
      </span>
    );
  }

  return (
    <span className="relative inline-block">
      <div
        className={`${sizeClasses[size]} bg-blue-600 rounded-full flex items-center justify-center font-medium text-white`}
      >
        {name.charAt(0).toUpperCase()}
      </div>
      {indicator}
    </span>
  );
}
