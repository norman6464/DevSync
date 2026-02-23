import { useState } from 'react';
import { Copy, Check } from 'lucide-react';

interface CodeBlockProps {
  code: string;
  language?: string;
  title?: string;
  showLineNumbers?: boolean;
  showCopy?: boolean;
  className?: string;
}

export default function CodeBlock({
  code,
  language,
  title,
  showLineNumbers = false,
  showCopy = false,
  className = '',
}: CodeBlockProps) {
  const [copied, setCopied] = useState(false);
  const lines = code.split('\n');

  const handleCopy = async () => {
    await navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className={`bg-gray-900 border border-gray-800 rounded-lg overflow-hidden ${className}`.trim()}>
      {(title || language || showCopy) && (
        <div className="flex items-center justify-between px-4 py-2 bg-gray-800/50 border-b border-gray-800">
          <div className="flex items-center gap-2">
            {title && <span className="text-xs text-gray-400">{title}</span>}
            {language && <span className="text-xs text-gray-500">{language}</span>}
          </div>
          {showCopy && (
            <button
              type="button"
              onClick={handleCopy}
              className="p-1 text-gray-400 hover:text-white transition-colors"
            >
              {copied ? <Check className="w-4 h-4 text-green-400" /> : <Copy className="w-4 h-4" />}
            </button>
          )}
        </div>
      )}
      <pre className="p-4 overflow-x-auto">
        <code className="text-sm font-mono text-gray-300">
          {showLineNumbers ? (
            lines.map((line, i) => (
              <div key={i} className="flex">
                <span className="select-none text-gray-600 w-8 text-right mr-4 flex-shrink-0">{i + 1}</span>
                <span>{line}</span>
              </div>
            ))
          ) : (
            code
          )}
        </code>
      </pre>
    </div>
  );
}
