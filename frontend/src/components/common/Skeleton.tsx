interface SkeletonProps {
  className?: string;
}

export function Skeleton({ className = '' }: SkeletonProps) {
  return (
    <div className={`animate-pulse bg-gray-800 rounded ${className}`} />
  );
}

export function PostCardSkeleton() {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-5">
      <div className="flex items-center gap-3 mb-4">
        <Skeleton className="w-10 h-10 rounded-full" />
        <div className="flex-1">
          <Skeleton className="h-4 w-24 mb-2" />
          <Skeleton className="h-3 w-16" />
        </div>
      </div>
      <Skeleton className="h-5 w-3/4 mb-3" />
      <Skeleton className="h-4 w-full mb-2" />
      <Skeleton className="h-4 w-5/6 mb-2" />
      <Skeleton className="h-4 w-2/3" />
      <div className="border-t border-gray-800 mt-4 pt-3 flex gap-4">
        <Skeleton className="h-4 w-16" />
        <Skeleton className="h-4 w-16" />
      </div>
    </div>
  );
}

export function UserCardSkeleton() {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
      <div className="flex items-center gap-3 mb-3">
        <Skeleton className="w-12 h-12 rounded-full" />
        <div className="flex-1">
          <Skeleton className="h-4 w-28 mb-2" />
          <Skeleton className="h-3 w-20" />
        </div>
      </div>
      <Skeleton className="h-3 w-full mb-1.5" />
      <Skeleton className="h-3 w-4/5" />
      <div className="border-t border-gray-800 mt-3 pt-3 flex gap-2">
        <Skeleton className="h-8 w-full rounded-lg" />
        <Skeleton className="h-8 w-full rounded-lg" />
      </div>
    </div>
  );
}

export function QuestionCardSkeleton() {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-5">
      <div className="flex items-start gap-3 mb-3">
        <div className="flex flex-col items-center gap-1">
          <Skeleton className="w-6 h-6" />
          <Skeleton className="w-8 h-4" />
        </div>
        <div className="flex-1">
          <Skeleton className="h-5 w-3/4 mb-2" />
          <Skeleton className="h-4 w-full mb-1.5" />
          <Skeleton className="h-4 w-5/6" />
          <div className="flex gap-2 mt-3">
            <Skeleton className="h-5 w-16 rounded-full" />
            <Skeleton className="h-5 w-20 rounded-full" />
          </div>
        </div>
      </div>
      <div className="flex items-center gap-3 text-xs">
        <Skeleton className="w-8 h-8 rounded-full" />
        <Skeleton className="h-3 w-24" />
        <Skeleton className="h-3 w-16" />
      </div>
    </div>
  );
}

export function ResourceCardSkeleton() {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-5">
      <div className="flex items-start justify-between mb-3">
        <Skeleton className="h-6 w-2/3" />
        <Skeleton className="h-6 w-12" />
      </div>
      <Skeleton className="h-4 w-full mb-2" />
      <Skeleton className="h-4 w-4/5 mb-4" />
      <div className="flex gap-2 mb-3">
        <Skeleton className="h-5 w-16 rounded-full" />
        <Skeleton className="h-5 w-20 rounded-full" />
        <Skeleton className="h-5 w-14 rounded-full" />
      </div>
      <div className="border-t border-gray-800 pt-3 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Skeleton className="w-6 h-6 rounded-full" />
          <Skeleton className="h-3 w-20" />
        </div>
        <div className="flex gap-3">
          <Skeleton className="h-4 w-12" />
          <Skeleton className="h-4 w-12" />
        </div>
      </div>
    </div>
  );
}

export function ListSkeleton({ count = 3 }: { count?: number }) {
  return (
    <>
      {Array.from({ length: count }).map((_, i) => (
        <Skeleton key={i} className="h-16 w-full mb-3" />
      ))}
    </>
  );
}

export function TableSkeleton({ rows = 5, cols = 4 }: { rows?: number; cols?: number }) {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md overflow-hidden">
      <div className="bg-gray-800 p-4 flex gap-4">
        {Array.from({ length: cols }).map((_, i) => (
          <Skeleton key={i} className="h-4 flex-1" />
        ))}
      </div>
      {Array.from({ length: rows }).map((_, i) => (
        <div key={i} className="p-4 flex gap-4 border-t border-gray-800">
          {Array.from({ length: cols }).map((_, j) => (
            <Skeleton key={j} className="h-4 flex-1" />
          ))}
        </div>
      ))}
    </div>
  );
}
