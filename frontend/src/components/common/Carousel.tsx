import { useState, ReactNode } from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';

interface CarouselProps<T> {
  items: T[];
  renderItem: (item: T, index: number) => ReactNode;
  loop?: boolean;
  showDots?: boolean;
  onChange?: (index: number) => void;
  className?: string;
}

export default function Carousel<T>({
  items,
  renderItem,
  loop = false,
  showDots = false,
  onChange,
  className = '',
}: CarouselProps<T>) {
  const [current, setCurrent] = useState(0);

  const goTo = (index: number) => {
    setCurrent(index);
    onChange?.(index);
  };

  const prev = () => {
    if (loop && current === 0) {
      goTo(items.length - 1);
    } else if (current > 0) {
      goTo(current - 1);
    }
  };

  const next = () => {
    if (loop && current === items.length - 1) {
      goTo(0);
    } else if (current < items.length - 1) {
      goTo(current + 1);
    }
  };

  const canPrev = loop || current > 0;
  const canNext = loop || current < items.length - 1;

  return (
    <div className={`relative ${className}`.trim()}>
      <div className="overflow-hidden">{renderItem(items[current], current)}</div>
      <button
        type="button"
        aria-label="前へ"
        onClick={prev}
        disabled={!canPrev}
        className="absolute left-2 top-1/2 -translate-y-1/2 p-1 bg-gray-800/80 rounded-full text-gray-300 hover:text-white disabled:opacity-30"
      >
        <ChevronLeft className="w-5 h-5" />
      </button>
      <button
        type="button"
        aria-label="次へ"
        onClick={next}
        disabled={!canNext}
        className="absolute right-2 top-1/2 -translate-y-1/2 p-1 bg-gray-800/80 rounded-full text-gray-300 hover:text-white disabled:opacity-30"
      >
        <ChevronRight className="w-5 h-5" />
      </button>
      {showDots && (
        <div className="flex justify-center gap-2 mt-3">
          {items.map((_, i) => (
            <button
              key={i}
              type="button"
              data-testid="carousel-dot"
              onClick={() => goTo(i)}
              className={`w-2 h-2 rounded-full transition-colors ${i === current ? 'bg-blue-500' : 'bg-gray-600'}`}
            />
          ))}
        </div>
      )}
    </div>
  );
}
