interface CheckboxOption {
  value: string;
  label: string;
  description?: string;
}

interface CheckboxGroupProps {
  options: CheckboxOption[];
  value: string[];
  onChange: (value: string[]) => void;
  label?: string;
  error?: string;
  disabled?: boolean;
  showSelectAll?: boolean;
  className?: string;
}

export default function CheckboxGroup({
  options,
  value,
  onChange,
  label,
  error,
  disabled = false,
  showSelectAll = false,
  className = '',
}: CheckboxGroupProps) {
  const handleToggle = (optionValue: string) => {
    if (value.includes(optionValue)) {
      onChange(value.filter((v) => v !== optionValue));
    } else {
      onChange([...value, optionValue]);
    }
  };

  const allSelected = options.every((o) => value.includes(o.value));

  return (
    <div className={`${className}`.trim()}>
      {label && <p className="text-sm text-gray-400 mb-2">{label}</p>}
      {showSelectAll && (
        <button
          type="button"
          onClick={() => onChange(allSelected ? [] : options.map((o) => o.value))}
          disabled={disabled}
          className="text-xs text-blue-400 hover:text-blue-300 mb-2 disabled:opacity-50"
        >
          {allSelected ? 'すべて解除' : 'すべて選択'}
        </button>
      )}
      <div className="space-y-2">
        {options.map((option) => (
          <label key={option.value} className="flex items-start gap-3 cursor-pointer">
            <input
              type="checkbox"
              checked={value.includes(option.value)}
              onChange={() => handleToggle(option.value)}
              disabled={disabled}
              className="mt-1 rounded border-gray-600 bg-gray-800 text-blue-500 focus:ring-blue-500"
            />
            <div>
              <span className="text-sm text-gray-200">{option.label}</span>
              {option.description && (
                <p className="text-xs text-gray-500">{option.description}</p>
              )}
            </div>
          </label>
        ))}
      </div>
      {error && <p className="mt-1 text-xs text-red-400">{error}</p>}
    </div>
  );
}
