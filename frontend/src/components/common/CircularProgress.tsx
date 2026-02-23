interface CircularProgressProps {
  value: number;
  size?: 'sm' | 'md' | 'lg';
  color?: string;
  showValue?: boolean;
  label?: string;
  className?: string;
}

const sizes = {
  sm: { width: 48, stroke: 4, fontSize: 'text-xs' },
  md: { width: 64, stroke: 5, fontSize: 'text-sm' },
  lg: { width: 96, stroke: 6, fontSize: 'text-lg' },
};

export default function CircularProgress({
  value,
  size = 'md',
  color = '#3b82f6',
  showValue = false,
  label,
  className = '',
}: CircularProgressProps) {
  const { width, stroke, fontSize } = sizes[size];
  const radius = (width - stroke) / 2;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (value / 100) * circumference;

  return (
    <div className={`inline-flex flex-col items-center gap-1 ${className}`.trim()}>
      <div className="relative inline-flex items-center justify-center">
        <svg width={width} height={width}>
          <circle
            cx={width / 2}
            cy={width / 2}
            r={radius}
            fill="none"
            stroke="#374151"
            strokeWidth={stroke}
          />
          <circle
            cx={width / 2}
            cy={width / 2}
            r={radius}
            fill="none"
            stroke={color}
            strokeWidth={stroke}
            strokeDasharray={circumference}
            strokeDashoffset={offset}
            strokeLinecap="round"
            transform={`rotate(-90 ${width / 2} ${width / 2})`}
            className="transition-all duration-500"
          />
        </svg>
        {showValue && (
          <span className={`absolute ${fontSize} font-semibold text-gray-200`}>
            {value}%
          </span>
        )}
      </div>
      {label && <span className="text-xs text-gray-400">{label}</span>}
    </div>
  );
}
