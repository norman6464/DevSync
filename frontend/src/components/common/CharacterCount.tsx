interface CharacterCountProps {
  current: number;
  max: number;
  showBar?: boolean;
  showRemaining?: boolean;
  className?: string;
}

export default function CharacterCount({
  current,
  max,
  showBar = false,
  showRemaining = false,
  className = '',
}: CharacterCountProps) {
  const percentage = (current / max) * 100;
  const isOver = current > max;
  const isWarning = percentage >= 80 && !isOver;

  const textColor = isOver ? 'text-red-400' : isWarning ? 'text-yellow-400' : 'text-gray-400';
  const barColor = isOver ? 'bg-red-400' : isWarning ? 'bg-yellow-400' : 'bg-blue-400';

  return (
    <div className={`${className}`.trim()}>
      <div className={`flex items-center gap-1 text-xs ${textColor}`}>
        {showRemaining ? (
          <>
            <span>{max - current}</span>
            <span>残り</span>
          </>
        ) : (
          <>
            <span>{current}</span>
            <span>/ {max}</span>
          </>
        )}
      </div>
      {showBar && (
        <div className="h-1 bg-gray-700 rounded-full mt-1 overflow-hidden">
          <div
            className={`h-full rounded-full transition-all ${barColor}`}
            style={{ width: `${Math.min(percentage, 100)}%` }}
          />
        </div>
      )}
    </div>
  );
}
