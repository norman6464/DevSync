const sizeClasses = {
  sm: 'text-sm',
  md: 'text-lg',
  lg: 'text-2xl',
};

interface PriceTagProps {
  price: number;
  originalPrice?: number;
  currency?: string;
  showDiscount?: boolean;
  size?: 'sm' | 'md' | 'lg';
  label?: string;
  className?: string;
}

function formatPrice(price: number, currency: string) {
  if (currency === '¥') return `¥${price.toLocaleString()}`;
  return `${currency}${price.toLocaleString()}`;
}

export default function PriceTag({
  price,
  originalPrice,
  currency = '¥',
  showDiscount = false,
  size = 'md',
  label,
  className = '',
}: PriceTagProps) {
  const discountPercent = originalPrice
    ? Math.round(((originalPrice - price) / originalPrice) * 100)
    : 0;

  return (
    <div className={`flex items-baseline gap-2 ${className}`.trim()}>
      <span className={`font-bold text-white ${sizeClasses[size]}`}>
        {price === 0 ? '無料' : formatPrice(price, currency)}
      </span>
      {originalPrice && originalPrice > price && (
        <span className="text-sm text-gray-500 line-through">
          {formatPrice(originalPrice, currency)}
        </span>
      )}
      {showDiscount && discountPercent > 0 && (
        <span className="text-sm text-green-400 font-medium">
          -{discountPercent}%
        </span>
      )}
      {label && <span className="text-xs text-gray-500">{label}</span>}
    </div>
  );
}
