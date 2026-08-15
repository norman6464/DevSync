import type { ReactNode } from 'react';

interface MarqueeProps {
  children: ReactNode;
  speed?: 'slow' | 'normal' | 'fast';
  direction?: 'left' | 'right';
  pauseOnHover?: boolean;
  className?: string;
}

export default function Marquee({
  children,
  speed = 'normal',
  direction = 'left',
  pauseOnHover = false,
  className = '',
}: MarqueeProps) {
  const animationClass = direction === 'right' ? 'animate-marquee-reverse' : 'animate-marquee';

  return (
    <div className={`overflow-hidden whitespace-nowrap ${className}`.trim()}>
      <div
        data-speed={speed}
        className={`inline-block ${animationClass} ${pauseOnHover ? 'hover:animation-paused' : ''}`}
      >
        {children}
      </div>
    </div>
  );
}
