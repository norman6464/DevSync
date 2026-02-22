import { ChevronLeft, ChevronRight } from 'lucide-react';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  showInfo?: boolean;
  maxVisible?: number;
}

export default function Pagination({
  currentPage,
  totalPages,
  onPageChange,
  showInfo = false,
  maxVisible = 7,
}: PaginationProps) {
  // ページが1つだけの場合は表示しない
  if (totalPages <= 1) {
    return null;
  }

  const handlePrevious = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1);
    }
  };

  const handleNext = () => {
    if (currentPage < totalPages) {
      onPageChange(currentPage + 1);
    }
  };

  const getPageNumbers = (): (number | string)[] => {
    const pages: (number | string)[] = [];

    if (totalPages <= maxVisible) {
      // 全ページを表示
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      // 省略記号を使用した表示
      const leftSiblingIndex = Math.max(currentPage - 1, 1);
      const rightSiblingIndex = Math.min(currentPage + 1, totalPages);

      const shouldShowLeftDots = leftSiblingIndex > 2;
      const shouldShowRightDots = rightSiblingIndex < totalPages - 1;

      if (!shouldShowLeftDots && shouldShowRightDots) {
        // 左側に省略なし、右側に省略あり
        const leftPages = Math.min(maxVisible - 2, totalPages - 2);
        for (let i = 1; i <= leftPages; i++) {
          pages.push(i);
        }
        pages.push('...');
        pages.push(totalPages);
      } else if (shouldShowLeftDots && !shouldShowRightDots) {
        // 左側に省略あり、右側に省略なし
        pages.push(1);
        pages.push('...');
        const rightPages = Math.min(maxVisible - 2, totalPages - 2);
        for (let i = totalPages - rightPages + 1; i <= totalPages; i++) {
          pages.push(i);
        }
      } else {
        // 両側に省略あり
        pages.push(1);
        pages.push('...');
        for (let i = leftSiblingIndex; i <= rightSiblingIndex; i++) {
          pages.push(i);
        }
        pages.push('...');
        pages.push(totalPages);
      }
    }

    return pages;
  };

  const pageNumbers = getPageNumbers();

  return (
    <nav className="flex items-center justify-center gap-2">
      {/* 前へボタン */}
      <button
        onClick={handlePrevious}
        disabled={currentPage === 1}
        className="p-2 text-gray-300 hover:text-white disabled:text-gray-600 disabled:cursor-not-allowed transition-colors rounded-md"
        aria-label="前のページ"
      >
        <ChevronLeft className="w-5 h-5" />
      </button>

      {/* ページ番号 */}
      <div className="flex items-center gap-1">
        {pageNumbers.map((page, index) => {
          if (page === '...') {
            return (
              <span key={`ellipsis-${index}`} className="px-3 py-1.5 text-gray-500">
                ...
              </span>
            );
          }

          const pageNumber = page as number;
          const isActive = pageNumber === currentPage;

          return (
            <button
              key={pageNumber}
              onClick={() => onPageChange(pageNumber)}
              className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${
                isActive
                  ? 'bg-blue-500 text-white'
                  : 'text-gray-300 hover:text-white hover:bg-gray-700'
              }`}
            >
              {pageNumber}
            </button>
          );
        })}
      </div>

      {/* 次へボタン */}
      <button
        onClick={handleNext}
        disabled={currentPage === totalPages}
        className="p-2 text-gray-300 hover:text-white disabled:text-gray-600 disabled:cursor-not-allowed transition-colors rounded-md"
        aria-label="次のページ"
      >
        <ChevronRight className="w-5 h-5" />
      </button>

      {/* ページ情報 */}
      {showInfo && (
        <div className="ml-4 text-sm text-gray-400">
          ページ {currentPage} / {totalPages}
        </div>
      )}
    </nav>
  );
}
