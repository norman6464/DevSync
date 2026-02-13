interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg';
  color?: 'primary' | 'secondary' | 'white';
  className?: string;
}

export default function LoadingSpinner({
  size = 'md',
  color = 'primary',
  className = ''
}: LoadingSpinnerProps) {
  const sizeClasses = {
    sm: 'w-5 h-5 border-2',
    md: 'w-8 h-8 border-4',
    lg: 'w-12 h-12 border-4',
  };

  const colorClasses = {
    primary: 'border-blue-500',
    secondary: 'border-gray-400',
    white: 'border-white',
  };

  return (
    <div className={`flex justify-center items-center py-12 ${className}`}>
      <div
        className={`${sizeClasses[size]} ${colorClasses[color]} border-t-transparent rounded-full animate-spin`}
        role="status"
        aria-label="読み込み中"
      />
    </div>
  );
}
