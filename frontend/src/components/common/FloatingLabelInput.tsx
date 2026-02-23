import { useState } from 'react';

interface FloatingLabelInputProps {
  value: string;
  onChange: (value: string) => void;
  label: string;
  type?: string;
  error?: string;
  disabled?: boolean;
  className?: string;
}

export default function FloatingLabelInput({
  value,
  onChange,
  label,
  type = 'text',
  error,
  disabled = false,
  className = '',
}: FloatingLabelInputProps) {
  const [focused, setFocused] = useState(false);
  const isFloating = focused || value.length > 0;

  return (
    <div className={`${className}`.trim()}>
      <div className="relative">
        <input
          type={type}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onFocus={() => setFocused(true)}
          onBlur={() => setFocused(false)}
          disabled={disabled}
          className={`w-full px-4 pt-5 pb-2 bg-gray-800 border rounded-lg text-gray-200 focus:outline-none transition-colors disabled:opacity-50 peer ${
            error ? 'border-red-500' : 'border-gray-700 focus:border-blue-500'
          }`}
        />
        <span
          className={`absolute left-4 transition-all pointer-events-none ${
            isFloating
              ? 'top-1 text-xs text-blue-400'
              : 'top-3.5 text-sm text-gray-500'
          }`}
        >
          {label}
        </span>
      </div>
      {error && <p className="mt-1 text-xs text-red-400">{error}</p>}
    </div>
  );
}
