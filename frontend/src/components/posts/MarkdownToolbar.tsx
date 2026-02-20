import { useTranslation } from 'react-i18next';
import { Link2, Image } from 'lucide-react';

interface MarkdownToolbarProps {
  onAction: (action: string) => void;
}

export default function MarkdownToolbar({ onAction }: MarkdownToolbarProps) {
  const { t } = useTranslation();

  const buttons = [
    { action: 'heading', icon: 'H', title: t('editor.heading') },
    { action: 'bold', icon: 'B', title: t('editor.bold'), className: 'font-bold' },
    { action: 'italic', icon: 'I', title: t('editor.italic'), className: 'italic' },
    { action: 'strikethrough', icon: 'S', title: t('editor.strikethrough'), className: 'line-through' },
    { action: 'divider' },
    { action: 'link', icon: <Link2 className="w-3.5 h-3.5" />, title: t('editor.link') },
    { action: 'image', icon: <Image className="w-3.5 h-3.5" />, title: t('editor.image') },
    { action: 'divider' },
    { action: 'code', icon: '<>', title: t('editor.inlineCode') },
    { action: 'codeblock', icon: '{}', title: t('editor.codeBlock') },
    { action: 'divider' },
    { action: 'quote', icon: '"', title: t('editor.quote') },
    { action: 'list', icon: '•', title: t('editor.bulletList') },
    { action: 'orderedlist', icon: '1.', title: t('editor.numberedList') },
    { action: 'task', icon: '☐', title: t('editor.taskList') },
  ];

  return (
    <div className="flex items-center gap-0.5 ml-auto px-2" role="toolbar" aria-label={t('editor.toolbar')}>
      {buttons.map((btn, i) =>
        btn.action === 'divider' ? (
          <div key={i} className="w-px h-4 bg-gray-600 mx-1" aria-hidden="true" />
        ) : (
          <button
            key={btn.action}
            type="button"
            onClick={() => onAction(btn.action)}
            className={`p-1.5 text-xs text-gray-400 hover:text-white hover:bg-gray-700 rounded transition-colors ${btn.className || ''}`}
            aria-label={btn.title}
            title={btn.title}
          >
            {btn.icon}
          </button>
        )
      )}
    </div>
  );
}
