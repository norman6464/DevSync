import { useState, useCallback, ReactNode, useRef } from 'react';

interface ResizablePanelProps {
  children: ReactNode;
  defaultSize?: number;
  minSize?: number;
  maxSize?: number;
  direction?: 'horizontal' | 'vertical';
  className?: string;
}

export default function ResizablePanel({
  children,
  defaultSize = 250,
  minSize = 100,
  maxSize = 800,
  direction = 'horizontal',
  className = '',
}: ResizablePanelProps) {
  const [size, setSize] = useState(defaultSize);
  const isResizing = useRef(false);

  const handleMouseDown = useCallback(
    (e: React.MouseEvent) => {
      e.preventDefault();
      isResizing.current = true;
      const startPos = direction === 'horizontal' ? e.clientX : e.clientY;
      const startSize = size;

      const handleMouseMove = (moveEvent: MouseEvent) => {
        if (!isResizing.current) return;
        const currentPos = direction === 'horizontal' ? moveEvent.clientX : moveEvent.clientY;
        const delta = currentPos - startPos;
        const newSize = Math.min(maxSize, Math.max(minSize, startSize + delta));
        setSize(newSize);
      };

      const handleMouseUp = () => {
        isResizing.current = false;
        document.removeEventListener('mousemove', handleMouseMove);
        document.removeEventListener('mouseup', handleMouseUp);
      };

      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
    },
    [size, minSize, maxSize, direction]
  );

  const isHorizontal = direction === 'horizontal';
  const style = isHorizontal ? { width: `${size}px` } : { height: `${size}px` };

  return (
    <div className={`relative flex ${isHorizontal ? 'flex-row' : 'flex-col'} ${className}`.trim()}>
      <div className="overflow-auto" style={style}>
        {children}
      </div>
      <div
        data-testid="resize-handle"
        onMouseDown={handleMouseDown}
        className={`flex-shrink-0 ${
          isHorizontal
            ? 'w-1 cursor-col-resize hover:bg-blue-500'
            : 'h-1 cursor-row-resize hover:bg-blue-500'
        } bg-gray-700 transition-colors`}
      />
    </div>
  );
}
