import { useState, useEffect } from 'react';
import { ArrowUp } from 'lucide-react';

interface ScrollToTopProps {
  threshold?: number;
  className?: string;
}

export default function ScrollToTop({
  threshold = 300,
  className = '',
}: ScrollToTopProps) {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setVisible(window.scrollY > threshold);
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, [threshold]);

  const scrollToTop = () => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  if (!visible) return null;

  return (
    <div className={`fixed bottom-6 right-6 z-40 ${className}`.trim()}>
      <button
        type="button"
        aria-label="トップへ戻る"
        onClick={scrollToTop}
        className="p-3 bg-blue-600 hover:bg-blue-500 text-white rounded-full shadow-lg transition-colors"
      >
        <ArrowUp className="w-5 h-5" />
      </button>
    </div>
  );
}
