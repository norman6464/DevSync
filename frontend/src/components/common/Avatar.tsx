import { sanitizeUrl } from '../../utils/url';

interface AvatarProps {
  name: string;
  avatarUrl?: string;
  size?: 'xs' | 'sm' | 'md' | 'lg';
}

const sizeClasses = {
  xs: 'w-6 h-6 text-xs',
  sm: 'w-8 h-8 text-sm',
  md: 'w-10 h-10 text-base',
  lg: 'w-16 h-16 text-xl',
};

export default function Avatar({ name, avatarUrl, size = 'md' }: AvatarProps) {
  const safeUrl = sanitizeUrl(avatarUrl);
  if (safeUrl) {
    return (
      <img
        src={safeUrl}
        alt={name}
        className={`${sizeClasses[size]} rounded-full object-cover`}
      />
    );
  }

  return (
    <div
      className={`${sizeClasses[size]} bg-blue-600 rounded-full flex items-center justify-center font-medium text-white`}
    >
      {name.charAt(0).toUpperCase()}
    </div>
  );
}
