interface RadioOption {
  value: string;
  label: string;
  description?: string;
}

interface RadioGroupProps {
  options: RadioOption[];
  value: string;
  onChange: (value: string) => void;
  label?: string;
  direction?: 'vertical' | 'horizontal';
  disabled?: boolean;
  className?: string;
}

export default function RadioGroup({
  options,
  value,
  onChange,
  label,
  direction = 'vertical',
  disabled = false,
  className = '',
}: RadioGroupProps) {
  const dirClass = direction === 'horizontal' ? 'flex-row' : 'flex-col';

  return (
    <div className={`${className}`.trim()}>
      {label && <label className="block text-sm text-gray-400 mb-2">{label}</label>}
      <div className={`flex gap-3 ${dirClass}`} role="radiogroup">
        {options.map((option) => (
          <label
            key={option.value}
            className={`flex items-start gap-3 cursor-pointer ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
          >
            <input
              type="radio"
              name="radio-group"
              value={option.value}
              checked={value === option.value}
              onChange={() => onChange(option.value)}
              disabled={disabled}
              className="mt-1 accent-blue-600"
            />
            <div>
              <span className="text-sm text-gray-200">{option.label}</span>
              {option.description && (
                <p className="text-xs text-gray-500 mt-0.5">{option.description}</p>
              )}
            </div>
          </label>
        ))}
      </div>
    </div>
  );
}
