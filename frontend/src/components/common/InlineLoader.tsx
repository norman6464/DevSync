interface InlineLoaderProps {
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export default function InlineLoader({
  size = 'md',
  className = '',
}: InlineLoaderProps) {
  const sizeClasses = {
    sm: 'w-4 h-4 border-2',
    md: 'w-5 h-5 border-2',
    lg: 'w-6 h-6 border-2',
  };

  return (
    <div className={`inline-flex items-center justify-center ${className}`}>
      <div
        className={`${sizeClasses[size]} border-blue-500 border-t-transparent rounded-full animate-spin`}
        role="status"
        aria-label="読み込み中"
      />
    </div>
  );
}
