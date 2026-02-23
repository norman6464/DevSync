import { useState, useRef, ReactNode, ReactElement } from 'react';

interface PopoverProps {
  trigger: ReactElement;
  children: ReactNode;
  position?: 'top' | 'bottom' | 'left' | 'right';
  title?: string;
  className?: string;
}

const positionClasses = {
  top: 'bottom-full left-1/2 -translate-x-1/2 mb-2',
  bottom: 'top-full left-1/2 -translate-x-1/2 mt-2',
  left: 'right-full top-1/2 -translate-y-1/2 mr-2',
  right: 'left-full top-1/2 -translate-y-1/2 ml-2',
};

export default function Popover({
  trigger,
  children,
  position = 'bottom',
  title,
  className = '',
}: PopoverProps) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  return (
    <div className="relative inline-block" ref={ref}>
      <div onClick={() => setOpen(!open)}>{trigger}</div>
      {open && (
        <div
          className={`absolute z-50 ${positionClasses[position]} bg-gray-800 border border-gray-700 rounded-lg shadow-xl p-3 min-w-[200px] ${className}`.trim()}
        >
          {title && <p className="text-sm font-semibold text-gray-200 mb-2">{title}</p>}
          {children}
        </div>
      )}
    </div>
  );
}
