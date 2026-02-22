const sizeConfig = {
  sm: { track: 'w-8 h-4', thumb: 'w-3 h-3', translate: 'translate-x-4' },
  md: { track: 'w-11 h-6', thumb: 'w-5 h-5', translate: 'translate-x-5' },
  lg: { track: 'w-14 h-7', thumb: 'w-6 h-6', translate: 'translate-x-7' },
};

interface ToggleProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  label?: string;
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  className?: string;
}

export default function Toggle({
  checked,
  onChange,
  label,
  size = 'md',
  disabled = false,
  className = '',
}: ToggleProps) {
  const config = sizeConfig[size];

  return (
    <div className={`flex items-center gap-3 ${className}`.trim()}>
      <button
        type="button"
        role="switch"
        aria-checked={checked}
        disabled={disabled}
        onClick={() => onChange(!checked)}
        className={`relative inline-flex flex-shrink-0 ${config.track} rounded-full transition-colors ${
          checked ? 'bg-blue-600' : 'bg-gray-600'
        } ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
      >
        <span
          className={`inline-block ${config.thumb} rounded-full bg-white shadow transform transition-transform ${
            checked ? config.translate : 'translate-x-0.5'
          } mt-0.5`}
        />
      </button>
      {label && (
        <span className="text-sm text-gray-300">{label}</span>
      )}
    </div>
  );
}
