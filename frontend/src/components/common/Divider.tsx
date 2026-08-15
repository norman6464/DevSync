import type { ReactNode } from 'react';

interface DividerProps {
  children?: ReactNode;
  orientation?: 'horizontal' | 'vertical';
  textAlign?: 'left' | 'center' | 'right';
  className?: string;
}

export default function Divider({
  children,
  orientation = 'horizontal',
  textAlign = 'center',
  className = '',
}: DividerProps) {
  // テキストがない場合はシンプルな線
  if (!children) {
    if (orientation === 'vertical') {
      return (
        <div
          className={`border-l border-gray-800 h-full ${className}`.trim()}
          role="separator"
          aria-orientation="vertical"
        />
      );
    }

    return (
      <div
        className={`border-t border-gray-800 w-full ${className}`.trim()}
        role="separator"
        aria-orientation="horizontal"
      />
    );
  }

  // テキスト付き区切り線（水平のみ）
  const alignClasses = {
    left: 'justify-start',
    center: 'justify-center',
    right: 'justify-end',
  };

  return (
    <div
      className={`flex items-center gap-4 ${alignClasses[textAlign]} w-full ${className}`.trim()}
      role="separator"
      aria-orientation="horizontal"
    >
      {(textAlign === 'center' || textAlign === 'right') && (
        <div className="flex-1 border-t border-gray-800" />
      )}
      <span className="text-sm text-gray-400">{children}</span>
      {(textAlign === 'center' || textAlign === 'left') && (
        <div className="flex-1 border-t border-gray-800" />
      )}
    </div>
  );
}
