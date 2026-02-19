import { Link } from 'react-router-dom';

interface MentionTextProps {
  text: string;
  className?: string;
}

const mentionRegex = /(?:^|(?<=\s))@([a-zA-Z0-9_-]+)/g;

export default function MentionText({ text, className = '' }: MentionTextProps) {
  const parts: (string | JSX.Element)[] = [];
  let lastIndex = 0;

  for (const match of text.matchAll(mentionRegex)) {
    const username = match[1];
    const start = match.index!;

    if (start > lastIndex) {
      parts.push(text.slice(lastIndex, start));
    }

    parts.push(
      <Link
        key={`${start}-${username}`}
        to={`/profile/${username}`}
        className="text-blue-400 hover:text-blue-300 font-medium transition-colors"
      >
        @{username}
      </Link>
    );

    lastIndex = start + match[0].length;
  }

  if (lastIndex < text.length) {
    parts.push(text.slice(lastIndex));
  }

  return <span className={className}>{parts}</span>;
}
