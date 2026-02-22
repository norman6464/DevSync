interface SkeletonProps {
  variant?: 'text' | 'circle' | 'rectangle';
  width?: string;
  height?: string;
  className?: string;
}

export default function Skeleton({
  variant = 'text',
  width,
  height,
  className = '',
}: SkeletonProps) {
  const variantClasses = {
    text: 'h-4 rounded',
    circle: 'rounded-full',
    rectangle: 'rounded',
  };

  const defaultSizes = {
    text: { width: '100%', height: undefined },
    circle: { width: '40px', height: '40px' },
    rectangle: { width: '100%', height: '100px' },
  };

  const finalWidth = width || defaultSizes[variant].width;
  const finalHeight = height || defaultSizes[variant].height;

  return (
    <div
      className={`bg-gray-800 animate-pulse ${variantClasses[variant]} ${className}`.trim()}
      style={{
        width: finalWidth,
        height: finalHeight,
      }}
    />
  );
}
