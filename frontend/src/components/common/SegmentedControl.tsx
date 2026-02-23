interface SegmentOption {
  value: string;
  label: string;
}

interface SegmentedControlProps {
  options: SegmentOption[];
  value: string;
  onChange: (value: string) => void;
  label?: string;
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  className?: string;
}

const sizeClasses = {
  sm: 'px-2 py-1 text-xs',
  md: 'px-3 py-1.5 text-sm',
  lg: 'px-4 py-2 text-base',
};

export default function SegmentedControl({
  options,
  value,
  onChange,
  label,
  size = 'md',
  disabled = false,
  className = '',
}: SegmentedControlProps) {
  return (
    <div className={`${className}`.trim()}>
      {label && <p className="text-sm text-gray-400 mb-1">{label}</p>}
      <div className="inline-flex bg-gray-800 rounded-lg p-1">
        {options.map((option) => (
          <button
            key={option.value}
            type="button"
            onClick={() => onChange(option.value)}
            disabled={disabled}
            className={`${sizeClasses[size]} rounded-md font-medium transition-colors disabled:opacity-50 ${
              value === option.value
                ? 'bg-blue-600 text-white'
                : 'text-gray-400 hover:text-gray-200'
            }`}
          >
            {option.label}
          </button>
        ))}
      </div>
    </div>
  );
}
