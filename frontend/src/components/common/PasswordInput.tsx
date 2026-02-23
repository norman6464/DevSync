import { useState } from 'react';
import { Eye, EyeOff } from 'lucide-react';

function getStrength(password: string): number {
  if (!password) return 0;
  let score = 0;
  if (password.length >= 8) score++;
  if (/[a-z]/.test(password) && /[A-Z]/.test(password)) score++;
  if (/\d/.test(password)) score++;
  if (/[^a-zA-Z0-9]/.test(password)) score++;
  return score;
}

const strengthColors = ['bg-red-500', 'bg-orange-500', 'bg-yellow-500', 'bg-green-500'];

interface PasswordInputProps {
  value: string;
  onChange: (value: string) => void;
  label?: string;
  error?: string;
  placeholder?: string;
  showStrength?: boolean;
  disabled?: boolean;
  className?: string;
}

export default function PasswordInput({
  value,
  onChange,
  label,
  error,
  placeholder = 'パスワード',
  showStrength = false,
  disabled = false,
  className = '',
}: PasswordInputProps) {
  const [visible, setVisible] = useState(false);
  const strength = getStrength(value);

  return (
    <div className={`${className}`.trim()}>
      {label && <label className="block text-sm text-gray-400 mb-1">{label}</label>}
      <div className="relative">
        <input
          type={visible ? 'text' : 'password'}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          disabled={disabled}
          className={`w-full pr-10 pl-4 py-2 bg-gray-800 border rounded-lg text-gray-200 placeholder-gray-500 focus:outline-none transition-colors disabled:opacity-50 ${
            error ? 'border-red-500' : 'border-gray-700 focus:border-blue-500'
          }`}
        />
        <button
          type="button"
          onClick={() => setVisible(!visible)}
          className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-white"
        >
          {visible ? <Eye className="w-4 h-4" /> : <EyeOff className="w-4 h-4" />}
        </button>
      </div>
      {showStrength && value && (
        <div className="flex gap-1 mt-2">
          {Array.from({ length: 4 }).map((_, i) => (
            <div
              key={i}
              data-testid="strength-bar"
              className={`h-1 flex-1 rounded-full ${i < strength ? strengthColors[strength - 1] : 'bg-gray-700'}`}
            />
          ))}
        </div>
      )}
      {error && <p className="mt-1 text-xs text-red-400">{error}</p>}
    </div>
  );
}
