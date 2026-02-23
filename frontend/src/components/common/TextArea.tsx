interface TextAreaProps {
  value: string;
  onChange: (value: string) => void;
  label?: string;
  error?: string;
  placeholder?: string;
  rows?: number;
  maxLength?: number;
  showCount?: boolean;
  disabled?: boolean;
  readOnly?: boolean;
  className?: string;
}

export default function TextArea({
  value,
  onChange,
  label,
  error,
  placeholder,
  rows = 3,
  maxLength,
  showCount = false,
  disabled = false,
  readOnly = false,
  className = '',
}: TextAreaProps) {
  return (
    <div className={`${className}`.trim()}>
      {label && <label className="block text-sm text-gray-400 mb-1">{label}</label>}
      <textarea
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        rows={rows}
        maxLength={maxLength}
        disabled={disabled}
        readOnly={readOnly}
        className={`w-full px-4 py-2 bg-gray-800 border rounded-lg text-gray-200 placeholder-gray-500 focus:outline-none transition-colors resize-y disabled:opacity-50 ${
          error ? 'border-red-500' : 'border-gray-700 focus:border-blue-500'
        }`}
      />
      <div className="flex justify-between mt-1">
        {error && <p className="text-xs text-red-400">{error}</p>}
        {showCount && maxLength && (
          <p className="text-xs text-gray-500 ml-auto">{value.length} / {maxLength}</p>
        )}
      </div>
    </div>
  );
}
