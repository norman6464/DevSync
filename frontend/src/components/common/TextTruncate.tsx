import { useState } from 'react';

interface TextTruncateProps {
  text: string;
  maxLength?: number;
  expandLabel?: string;
  collapseLabel?: string;
  className?: string;
}

export default function TextTruncate({
  text,
  maxLength = 200,
  expandLabel = 'もっと見る',
  collapseLabel = '閉じる',
  className = '',
}: TextTruncateProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const shouldTruncate = text.length > maxLength;

  return (
    <div className={`${className}`.trim()}>
      <span className="text-gray-300">
        {isExpanded || !shouldTruncate ? text : text.slice(0, maxLength)}
        {!isExpanded && shouldTruncate && <span>...</span>}
      </span>
      {shouldTruncate && (
        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="ml-1 text-sm text-blue-400 hover:text-blue-300 transition-colors"
        >
          {isExpanded ? collapseLabel : expandLabel}
        </button>
      )}
    </div>
  );
}
