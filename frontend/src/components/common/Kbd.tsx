const sizeClasses = {
  sm: 'text-xs px-1 py-0.5',
  md: 'text-sm px-1.5 py-0.5',
  lg: 'text-base px-2 py-1',
};

interface KbdProps {
  keys: string[];
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export default function Kbd({
  keys,
  size = 'md',
  className = '',
}: KbdProps) {
  return (
    <span className={`inline-flex items-center gap-1 ${className}`.trim()}>
      {keys.map((key, index) => (
        <span key={index} className="inline-flex items-center gap-1">
          {index > 0 && <span className="text-gray-500">+</span>}
          <kbd
            className={`${sizeClasses[size]} rounded border border-gray-600 bg-gray-800 text-gray-300 font-mono shadow-sm`}
          >
            {key}
          </kbd>
        </span>
      ))}
    </span>
  );
}
