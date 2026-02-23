import { useState } from 'react';
import { X, ChevronLeft, ChevronRight } from 'lucide-react';

interface ImageItem {
  src: string;
  alt: string;
}

interface ImageGalleryProps {
  images: ImageItem[];
  columns?: number;
  emptyMessage?: string;
  className?: string;
}

export default function ImageGallery({
  images,
  columns = 3,
  emptyMessage = '画像がありません',
  className = '',
}: ImageGalleryProps) {
  const [lightboxIndex, setLightboxIndex] = useState<number | null>(null);

  if (images.length === 0) {
    return (
      <div className={`${className}`.trim()}>
        <p className="text-center text-gray-500 py-8">{emptyMessage}</p>
      </div>
    );
  }

  const closeLightbox = () => setLightboxIndex(null);
  const prev = () => setLightboxIndex((i) => (i != null && i > 0 ? i - 1 : i));
  const next = () => setLightboxIndex((i) => (i != null && i < images.length - 1 ? i + 1 : i));

  return (
    <div className={`${className}`.trim()}>
      <div
        className="grid gap-2"
        style={{ gridTemplateColumns: `repeat(${columns}, minmax(0, 1fr))` }}
      >
        {images.map((img, i) => (
          <button
            key={i}
            type="button"
            onClick={() => setLightboxIndex(i)}
            className="overflow-hidden rounded-lg hover:opacity-80 transition-opacity"
          >
            <img src={img.src} alt={img.alt} className="w-full h-full object-cover" />
          </button>
        ))}
      </div>

      {lightboxIndex != null && (
        <div
          data-testid="lightbox"
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/80"
        >
          <button
            type="button"
            aria-label="閉じる"
            onClick={closeLightbox}
            className="absolute top-4 right-4 text-white hover:text-gray-300"
          >
            <X className="w-6 h-6" />
          </button>
          <button
            type="button"
            aria-label="前へ"
            onClick={prev}
            disabled={lightboxIndex === 0}
            className="absolute left-4 text-white hover:text-gray-300 disabled:opacity-30"
          >
            <ChevronLeft className="w-8 h-8" />
          </button>
          <img
            src={images[lightboxIndex].src}
            alt={images[lightboxIndex].alt}
            className="max-h-[80vh] max-w-[80vw] object-contain"
          />
          <button
            type="button"
            aria-label="次へ"
            onClick={next}
            disabled={lightboxIndex === images.length - 1}
            className="absolute right-4 text-white hover:text-gray-300 disabled:opacity-30"
          >
            <ChevronRight className="w-8 h-8" />
          </button>
        </div>
      )}
    </div>
  );
}
