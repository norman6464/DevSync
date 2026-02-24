import { ReactNode } from 'react';

interface WatermarkTextProps {
  text: string;
  children: ReactNode;
  opacity?: number;
  rotate?: number;
  fontSize?: string;
  className?: string;
}

export default function WatermarkText({
  text,
  children,
  opacity = 10,
  rotate = -45,
  fontSize = '3rem',
  className = '',
}: WatermarkTextProps) {
  return (
    <div className={`relative overflow-hidden ${className}`.trim()}>
      <span
        className="absolute inset-0 flex items-center justify-center pointer-events-none select-none text-gray-400 font-bold opacity-10"
        style={{
          opacity: opacity / 100,
          transform: `rotate(${rotate}deg)`,
          fontSize,
        }}
      >
        {text}
      </span>
      <div className="relative z-10">{children}</div>
    </div>
  );
}
