import { ReactNode, Children } from 'react';

interface MasonryGridProps {
  children: ReactNode;
  columns?: number;
  gap?: number;
  className?: string;
}

export default function MasonryGrid({
  children,
  columns = 3,
  gap = 8,
  className = '',
}: MasonryGridProps) {
  const items = Children.toArray(children);

  const columnItems: ReactNode[][] = Array.from({ length: columns }, () => []);
  items.forEach((item, i) => {
    columnItems[i % columns].push(item);
  });

  return (
    <div className={`flex ${className}`.trim()} style={{ gap: `${gap}px` }}>
      {columnItems.map((colItems, i) => (
        <div
          key={i}
          data-testid="masonry-column"
          className="flex-1 flex flex-col"
          style={{ gap: `${gap}px` }}
        >
          {colItems}
        </div>
      ))}
    </div>
  );
}
