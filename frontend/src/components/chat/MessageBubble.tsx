import { format } from 'date-fns';
import Avatar from '../common/Avatar';

interface MessageBubbleProps {
  content: string;
  createdAt: string;
  isOwn: boolean;
  senderName?: string;
  senderAvatarUrl?: string;
  showSenderInfo?: boolean;
}

export default function MessageBubble({
  content,
  createdAt,
  isOwn,
  senderName,
  senderAvatarUrl,
  showSenderInfo = false,
}: MessageBubbleProps) {
  const time = format(new Date(createdAt), 'HH:mm');

  return (
    <div className={`flex items-end gap-2 ${isOwn ? 'justify-end' : 'justify-start'}`}>
      {!isOwn && showSenderInfo && (
        <Avatar name={senderName || ''} avatarUrl={senderAvatarUrl} size="xs" />
      )}
      {isOwn && (
        <span className="chat-bubble-time-own text-xs mb-0.5">{time}</span>
      )}
      <div>
        {!isOwn && showSenderInfo && (
          <p className="text-xs text-gray-400 mb-1">{senderName}</p>
        )}
        <div
          className={`px-4 py-2.5 rounded-md ${
            isOwn ? 'chat-bubble-own rounded-br-md' : 'chat-bubble-other rounded-bl-md'
          }`}
        >
          <p className="text-sm leading-relaxed">{content}</p>
        </div>
      </div>
      {!isOwn && (
        <span className="chat-bubble-time-other text-xs mb-0.5">{time}</span>
      )}
    </div>
  );
}
