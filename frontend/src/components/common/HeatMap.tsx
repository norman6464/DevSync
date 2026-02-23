interface HeatMapData {
  date: string;
  count: number;
}

const cellSizeMap = {
  sm: 'w-3 h-3',
  md: 'w-4 h-4',
  lg: 'w-5 h-5',
};

function getIntensityClass(count: number): string {
  if (count === 0) return 'bg-gray-800';
  if (count <= 3) return 'bg-green-900';
  if (count <= 6) return 'bg-green-700';
  return 'bg-green-500';
}

interface HeatMapProps {
  data: HeatMapData[];
  cellSize?: 'sm' | 'md' | 'lg';
  showLegend?: boolean;
  className?: string;
}

export default function HeatMap({
  data,
  cellSize = 'md',
  showLegend = false,
  className = '',
}: HeatMapProps) {
  return (
    <div className={`${className}`.trim()}>
      <div className="flex flex-wrap gap-1">
        {data.map((item) => (
          <div
            key={item.date}
            data-testid="heatmap-cell"
            data-date={item.date}
            data-count={String(item.count)}
            className={`${cellSizeMap[cellSize]} rounded-sm ${getIntensityClass(item.count)}`}
            title={`${item.date}: ${item.count}`}
          />
        ))}
      </div>
      {showLegend && (
        <div className="flex items-center gap-1 mt-2 text-xs text-gray-500">
          <span>少</span>
          <div className="w-3 h-3 rounded-sm bg-gray-800" />
          <div className="w-3 h-3 rounded-sm bg-green-900" />
          <div className="w-3 h-3 rounded-sm bg-green-700" />
          <div className="w-3 h-3 rounded-sm bg-green-500" />
          <span>多</span>
        </div>
      )}
    </div>
  );
}
