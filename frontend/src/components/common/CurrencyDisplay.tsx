import { TrendingUp, TrendingDown } from 'lucide-react';

interface CurrencyDisplayProps {
  amount: number;
  currency: 'JPY' | 'USD' | 'EUR';
  change?: number;
  label?: string;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

const currencySymbols = { JPY: '¥', USD: '$', EUR: '€' };
const sizeClasses = { sm: 'text-lg', md: 'text-2xl', lg: 'text-3xl' };

function formatAmount(amount: number, currency: string): string {
  const formatted = amount.toLocaleString('ja-JP');
  return `${currencySymbols[currency as keyof typeof currencySymbols]}${formatted}`;
}

export default function CurrencyDisplay({
  amount,
  currency,
  change,
  label,
  size = 'md',
  className = '',
}: CurrencyDisplayProps) {
  return (
    <div className={`${className}`.trim()}>
      {label && <p className="text-xs text-gray-500 mb-1">{label}</p>}
      <p className={`${sizeClasses[size]} font-bold text-gray-200`}>
        {formatAmount(amount, currency)}
      </p>
      {change != null && (
        <div className={`flex items-center gap-1 mt-1 text-sm ${change >= 0 ? 'text-green-400' : 'text-red-400'}`}>
          {change >= 0 ? (
            <TrendingUp className="w-4 h-4" />
          ) : (
            <TrendingDown className="w-4 h-4" />
          )}
          <span>{change >= 0 ? '+' : ''}{change}%</span>
        </div>
      )}
    </div>
  );
}
