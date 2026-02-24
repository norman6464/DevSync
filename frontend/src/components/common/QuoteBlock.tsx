import { Quote } from 'lucide-react';

interface QuoteBlockProps {
  quote: string;
  author?: string;
  source?: string;
  color?: 'gray' | 'blue' | 'green' | 'purple';
  className?: string;
}

const colorClasses = {
  gray: 'border-gray-500 text-gray-400',
  blue: 'border-blue-500 text-blue-400',
  green: 'border-green-500 text-green-400',
  purple: 'border-purple-500 text-purple-400',
};

export default function QuoteBlock({
  quote,
  author,
  source,
  color = 'gray',
  className = '',
}: QuoteBlockProps) {
  return (
    <blockquote
      className={`border-l-4 pl-4 py-2 ${colorClasses[color]} ${className}`.trim()}
    >
      <Quote className="w-5 h-5 mb-2 opacity-50" />
      <p className="text-gray-200 text-sm italic leading-relaxed">{quote}</p>
      {(author || source) && (
        <footer className="mt-2 text-xs">
          {author && <cite className="not-italic font-medium">{author}</cite>}
          {author && source && <span className="mx-1">—</span>}
          {source && <span className="opacity-70">{source}</span>}
        </footer>
      )}
    </blockquote>
  );
}
