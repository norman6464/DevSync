interface SliderInputProps {
  value: number;
  onChange: (value: number) => void;
  min?: number;
  max?: number;
  step?: number;
  label?: string;
  showValue?: boolean;
  suffix?: string;
  disabled?: boolean;
  className?: string;
}

export default function SliderInput({
  value,
  onChange,
  min = 0,
  max = 100,
  step = 1,
  label,
  showValue = false,
  suffix,
  disabled = false,
  className = '',
}: SliderInputProps) {
  return (
    <div className={`${className}`.trim()}>
      {(label || showValue) && (
        <div className="flex items-center justify-between mb-2">
          {label && <label className="text-sm text-gray-400">{label}</label>}
          {showValue && (
            <span className="text-sm text-gray-200 font-medium">
              {value}{suffix && <span className="text-gray-400 ml-0.5">{suffix}</span>}
            </span>
          )}
        </div>
      )}
      <input
        type="range"
        value={value}
        onChange={(e) => onChange(Number(e.target.value))}
        min={min}
        max={max}
        step={step}
        disabled={disabled}
        className="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer accent-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
      />
    </div>
  );
}
