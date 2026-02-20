import { useState, useEffect } from 'react';
import { ChevronUp } from 'lucide-react';
import { useTranslation } from 'react-i18next';

export default function ScrollToTop() {
  const { t } = useTranslation();
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const onScroll = () => setVisible(window.scrollY > 200);
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  }, []);

  if (!visible) return null;

  return (
    <button
      onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
      className="fixed bottom-6 right-6 z-40 p-3 bg-gray-800 hover:bg-gray-700 text-gray-300 hover:text-white border border-gray-700 rounded-full shadow-sm transition-all duration-200"
      aria-label={t('common.scrollToTop')}
    >
      <ChevronUp className="w-5 h-5" />
    </button>
  );
}
