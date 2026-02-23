import { useEffect, useRef, ReactNode } from 'react';
import { Loader2 } from 'lucide-react';

interface InfiniteScrollProps {
  children: ReactNode;
  onLoadMore: () => void;
  hasMore: boolean;
  loading?: boolean;
  loadingText?: string;
  endMessage?: string;
  className?: string;
}

export default function InfiniteScroll({
  children,
  onLoadMore,
  hasMore,
  loading = false,
  loadingText,
  endMessage,
  className = '',
}: InfiniteScrollProps) {
  const sentinelRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!hasMore || loading) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          onLoadMore();
        }
      },
      { threshold: 0.1 }
    );

    const sentinel = sentinelRef.current;
    if (sentinel) {
      observer.observe(sentinel);
    }

    return () => observer.disconnect();
  }, [hasMore, loading, onLoadMore]);

  return (
    <div className={`${className}`.trim()}>
      {children}
      {hasMore && (
        <div ref={sentinelRef} data-testid="sentinel" className="h-1" />
      )}
      {loading && (
        <div className="flex items-center justify-center gap-2 py-4">
          <Loader2 className="w-5 h-5 text-gray-400 animate-spin" />
          {loadingText && <span className="text-sm text-gray-400">{loadingText}</span>}
        </div>
      )}
      {!hasMore && endMessage && (
        <div className="text-center py-4 text-sm text-gray-500">
          {endMessage}
        </div>
      )}
    </div>
  );
}
