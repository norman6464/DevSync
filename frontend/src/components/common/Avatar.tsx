import { User } from 'lucide-react';

export type AvatarSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl';

export interface AvatarProps {
  /** 画像 URL。ProfileCard と揃えた名前で受ける。 */
  avatarUrl?: string;
  alt?: string;
  name?: string;
  size?: AvatarSize;
  online?: boolean;
  rounded?: boolean;
  className?: string;
}

const sizeClasses: Record<AvatarSize, string> = {
  xs: 'w-6 h-6 text-[10px]',
  sm: 'w-8 h-8 text-xs',
  md: 'w-10 h-10 text-sm',
  lg: 'w-12 h-12 text-base',
  xl: 'w-16 h-16 text-lg',
};

const iconSizes: Record<AvatarSize, string> = {
  xs: 'w-3 h-3',
  sm: 'w-4 h-4',
  md: 'w-5 h-5',
  lg: 'w-6 h-6',
  xl: 'w-8 h-8',
};

export default function Avatar({
  avatarUrl,
  alt,
  name,
  size = 'md',
  online,
  rounded = true,
  className = '',
}: AvatarProps) {
  const shape = rounded ? 'rounded-full' : 'rounded-lg';
  const initial = name ? name.charAt(0) : null;

  return (
    <div className={`relative inline-flex ${className}`.trim()}>
      <div
        className={`${sizeClasses[size]} ${shape} bg-gray-700 flex items-center justify-center overflow-hidden`}
      >
        {avatarUrl ? (
          <img src={avatarUrl} alt={alt || name || ''} className="w-full h-full object-cover" />
        ) : initial ? (
          <span className="font-medium text-gray-200">{initial}</span>
        ) : (
          <User className={`${iconSizes[size]} text-gray-400`} />
        )}
      </div>
      {online != null && (
        <span
          className={`absolute bottom-0 right-0 w-3 h-3 border-2 border-gray-900 rounded-full ${
            online ? 'bg-green-500' : 'bg-gray-500'
          }`}
        />
      )}
    </div>
  );
}
