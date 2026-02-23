import { Plus, Minus } from 'lucide-react';

interface NumberInputProps {
  value: number;
  onChange: (value: number) => void;
  min?: number;
  max?: number;
  step?: number;
  disabled?: boolean;
  className?: string;
}

export default function NumberInput({
  value,
  onChange,
  min,
  max,
  step = 1,
  disabled = false,
  className = '',
}: NumberInputProps) {
  const canDecrement = min === undefined || value - step >= min;
  const canIncrement = max === undefined || value + step <= max;

  const handleIncrement = () => {
    if (canIncrement && !disabled) onChange(value + step);
  };

  const handleDecrement = () => {
    if (canDecrement && !disabled) onChange(value - step);
  };

  return (
    <div className={`flex items-center ${className}`.trim()}>
      <button
        type="button"
        aria-label="減少"
        onClick={handleDecrement}
        disabled={disabled || !canDecrement}
        className="px-3 py-2 bg-gray-800 border border-gray-700 rounded-l-lg text-gray-400 hover:text-white hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        <Minus className="w-4 h-4" />
      </button>
      <input
        type="number"
        value={value}
        onChange={(e) => onChange(Number(e.target.value))}
        disabled={disabled}
        className="w-20 text-center py-2 bg-gray-800 border-y border-gray-700 text-gray-200 focus:outline-none disabled:opacity-50 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
      />
      <button
        type="button"
        aria-label="増加"
        onClick={handleIncrement}
        disabled={disabled || !canIncrement}
        className="px-3 py-2 bg-gray-800 border border-gray-700 rounded-r-lg text-gray-400 hover:text-white hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        <Plus className="w-4 h-4" />
      </button>
    </div>
  );
}
