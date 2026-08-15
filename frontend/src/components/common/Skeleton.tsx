interface SkeletonProps {
  variant?: 'text' | 'circle' | 'rectangle';
  width?: string;
  height?: string;
  className?: string;
}

export function Skeleton({
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

export default Skeleton;

/** カード 1 枚分の骨組み。各ページのローディング表示で使う。 */
function CardFrame({ children, className = '' }: { children: React.ReactNode; className?: string }) {
  return <div className={`bg-gray-900 border border-gray-800 rounded-md ${className}`.trim()}>{children}</div>;
}

/** 投稿カードのローディング表示。 */
export function PostCardSkeleton() {
  return (
    <CardFrame className="p-5">
      <div className="flex items-center gap-3 mb-4">
        <Skeleton variant="circle" width="40px" height="40px" />
        <div className="flex-1">
          <Skeleton width="6rem" height="1rem" className="mb-2" />
          <Skeleton width="4rem" height="0.75rem" />
        </div>
      </div>
      <Skeleton width="75%" height="1.25rem" className="mb-3" />
      <Skeleton height="1rem" className="mb-2" />
      <Skeleton width="83%" height="1rem" className="mb-2" />
      <Skeleton width="66%" height="1rem" />
      <div className="border-t border-gray-800 mt-4 pt-3 flex gap-4">
        <Skeleton width="4rem" height="1rem" />
        <Skeleton width="4rem" height="1rem" />
      </div>
    </CardFrame>
  );
}

/** ユーザーカードのローディング表示。 */
export function UserCardSkeleton() {
  return (
    <CardFrame className="p-4">
      <div className="flex items-center gap-3 mb-3">
        <Skeleton variant="circle" width="48px" height="48px" />
        <div className="flex-1">
          <Skeleton width="7rem" height="1rem" className="mb-2" />
          <Skeleton width="5rem" height="0.75rem" />
        </div>
      </div>
      <Skeleton height="0.75rem" className="mb-1.5" />
      <Skeleton width="80%" height="0.75rem" />
      <div className="border-t border-gray-800 mt-3 pt-3 flex gap-2">
        <Skeleton variant="rectangle" height="2rem" className="rounded-lg" />
        <Skeleton variant="rectangle" height="2rem" className="rounded-lg" />
      </div>
    </CardFrame>
  );
}

/** 質問カードのローディング表示。 */
export function QuestionCardSkeleton() {
  return (
    <CardFrame className="p-5">
      <Skeleton width="70%" height="1.25rem" className="mb-3" />
      <Skeleton height="1rem" className="mb-2" />
      <Skeleton width="60%" height="1rem" className="mb-4" />
      <div className="flex gap-2">
        <Skeleton width="3.5rem" height="1.5rem" className="rounded-full" />
        <Skeleton width="3.5rem" height="1.5rem" className="rounded-full" />
      </div>
    </CardFrame>
  );
}

/** 学習リソースカードのローディング表示。 */
export function ResourceCardSkeleton() {
  return (
    <CardFrame className="p-5">
      <div className="flex items-start gap-3 mb-3">
        <Skeleton variant="rectangle" width="56px" height="56px" className="rounded-lg" />
        <div className="flex-1">
          <Skeleton width="65%" height="1.25rem" className="mb-2" />
          <Skeleton width="40%" height="0.75rem" />
        </div>
      </div>
      <Skeleton height="1rem" className="mb-2" />
      <Skeleton width="85%" height="1rem" />
    </CardFrame>
  );
}
