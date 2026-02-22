import { useState } from 'react';
import { Star } from 'lucide-react';

const sizeMap = {
  sm: 'w-4 h-4',
  md: 'w-6 h-6',
  lg: 'w-8 h-8',
};

interface RatingProps {
  value: number;
  onChange?: (value: number) => void;
  max?: number;
  readOnly?: boolean;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export default function Rating({
  value,
  onChange,
  max = 5,
  readOnly = false,
  size = 'md',
  className = '',
}: RatingProps) {
  const [hoverValue, setHoverValue] = useState<number | null>(null);

  const displayValue = hoverValue ?? value;

  const handleClick = (index: number) => {
    if (readOnly || !onChange) return;
    onChange(index === value ? 0 : index);
  };

  const stars = Array.from({ length: max }, (_, i) => i + 1);

  return (
    <div className={`flex items-center gap-1 ${className}`.trim()}>
      {stars.map((starValue) => {
        const isFilled = starValue <= displayValue;
        const starClasses = `${sizeMap[size]} ${isFilled ? 'text-yellow-400 fill-yellow-400' : 'text-gray-600'}`;

        if (readOnly) {
          return (
            <Star key={starValue} className={starClasses} />
          );
        }

        return (
          <button
            key={starValue}
            type="button"
            onClick={() => handleClick(starValue)}
            onMouseEnter={() => setHoverValue(starValue)}
            onMouseLeave={() => setHoverValue(null)}
            className="p-0 bg-transparent border-none cursor-pointer"
          >
            <Star className={starClasses} />
          </button>
        );
      })}
    </div>
  );
}
