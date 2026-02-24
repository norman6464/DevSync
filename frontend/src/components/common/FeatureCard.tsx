interface FeatureCardProps {
  title: string;
  description: string;
  icon?: string;
  badge?: string;
  onClick?: () => void;
  disabled?: boolean;
  className?: string;
}

export default function FeatureCard({
  title,
  description,
  icon,
  badge,
  onClick,
  disabled = false,
  className = '',
}: FeatureCardProps) {
  return (
    <div
      onClick={disabled ? undefined : onClick}
      className={`p-4 bg-gray-800/50 border border-gray-700 rounded-xl hover:border-gray-600 transition-colors ${
        onClick && !disabled ? 'cursor-pointer' : ''
      } ${disabled ? 'opacity-50' : ''} ${className}`.trim()}
    >
      <div className="flex items-start justify-between">
        <div className="flex items-start gap-3">
          {icon && <span className="text-2xl">{icon}</span>}
          <div>
            <h3 className="text-sm font-semibold text-gray-200">{title}</h3>
            <p className="mt-1 text-xs text-gray-400 leading-relaxed">{description}</p>
          </div>
        </div>
        {badge && (
          <span className="px-2 py-0.5 text-xs font-medium bg-blue-600/20 text-blue-400 rounded-full">
            {badge}
          </span>
        )}
      </div>
    </div>
  );
}
