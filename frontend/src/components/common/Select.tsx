import { useState, useRef, useEffect } from 'react';
import { ChevronDown } from 'lucide-react';

interface SelectOption {
  value: string;
  label: string;
}

interface SelectProps {
  options: SelectOption[];
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  label?: string;
  error?: string;
  disabled?: boolean;
  className?: string;
}

export default function Select({
  options,
  value,
  onChange,
  placeholder = '',
  label,
  error,
  disabled = false,
  className = '',
}: SelectProps) {
  const [isOpen, setIsOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  const selectedOption = options.find((o) => o.value === value);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleSelect = (optValue: string) => {
    onChange(optValue);
    setIsOpen(false);
  };

  return (
    <div ref={ref} className={`${className}`.trim()}>
      {label && <label className="block text-sm text-gray-400 mb-1">{label}</label>}
      <button
        type="button"
        onClick={() => !disabled && setIsOpen(!isOpen)}
        className={`w-full flex items-center justify-between px-4 py-2 bg-gray-800 border rounded-lg text-left transition-colors ${
          error ? 'border-red-500' : 'border-gray-700 focus:border-blue-500'
        } ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
      >
        <span className={selectedOption ? 'text-gray-200' : 'text-gray-500'}>
          {selectedOption?.label || placeholder}
        </span>
        <ChevronDown className={`w-4 h-4 text-gray-400 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
      </button>
      {isOpen && (
        <div className="absolute z-50 mt-1 w-full bg-gray-800 border border-gray-700 rounded-lg shadow-lg max-h-60 overflow-y-auto">
          {options.map((option) => (
            <button
              key={option.value}
              type="button"
              onClick={() => handleSelect(option.value)}
              className={`w-full text-left px-4 py-2 text-sm transition-colors ${
                option.value === value ? 'bg-blue-600/20 text-blue-400' : 'text-gray-200 hover:bg-gray-700'
              }`}
            >
              {option.label}
            </button>
          ))}
        </div>
      )}
      {error && <p className="mt-1 text-xs text-red-400">{error}</p>}
    </div>
  );
}
