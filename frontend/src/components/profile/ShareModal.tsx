import { useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import html2canvas from 'html2canvas';
import toast from 'react-hot-toast';
import ShareableProfileCard from './ShareableProfileCard';
import ShareButtons from './ShareButtons';
import type { User } from '../../types/user';
import type { GitHubLanguageStat } from '../../types/github';

interface ShareModalProps {
  isOpen: boolean;
  onClose: () => void;
  user: User;
  followerCount: number;
  followingCount: number;
  totalContributions: number;
  languages: GitHubLanguageStat[];
  postCount: number;
}

export default function ShareModal({
  isOpen,
  onClose,
  user,
  followerCount,
  followingCount,
  totalContributions,
  languages,
  postCount,
}: ShareModalProps) {
  const { t } = useTranslation();
  const cardRef = useRef<HTMLDivElement>(null);
  const [generating, setGenerating] = useState(false);

  if (!isOpen) return null;

  const profileUrl = `${window.location.origin}/profile/${encodeURIComponent(user.username)}`;

  const generateImage = async (): Promise<Blob | null> => {
    if (!cardRef.current) return null;

    try {
      const canvas = await html2canvas(cardRef.current, {
        backgroundColor: null,
        scale: 2,
        useCORS: true,
        logging: false,
      });

      return new Promise((resolve) => {
        canvas.toBlob((blob) => resolve(blob), 'image/png', 1.0);
      });
    } catch (error) {
      console.error('Failed to generate image:', error);
      return null;
    }
  };

  const handleDownload = async () => {
    setGenerating(true);
    try {
      const blob = await generateImage();
      if (!blob) {
        toast.error(t('sharing.downloadFailed'));
        return;
      }

      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `devsync-${(user.github_username || user.name || 'profile').replace(/[^a-zA-Z0-9_-]/g, '_')}.png`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
      toast.success(t('sharing.downloaded'));
    } finally {
      setGenerating(false);
    }
  };

  const handleCopyLink = async () => {
    try {
      await navigator.clipboard.writeText(profileUrl);
      toast.success(t('sharing.linkCopied'));
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  };

  const handleShareTwitter = () => {
    const text = t('sharing.twitterText', { name: user.name || 'Developer' });
    const url = `https://twitter.com/intent/tweet?text=${encodeURIComponent(text)}&url=${encodeURIComponent(profileUrl)}`;
    window.open(url, '_blank', 'width=600,height=400,noopener,noreferrer');
  };

  const handleShareLinkedIn = () => {
    const url = `https://www.linkedin.com/sharing/share-offsite/?url=${encodeURIComponent(profileUrl)}`;
    window.open(url, '_blank', 'width=600,height=400,noopener,noreferrer');
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/70 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Modal */}
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="share-modal-title"
        className="relative bg-gray-900 border border-gray-800 rounded-md shadow-sm max-w-3xl w-full mx-4 max-h-[90vh] overflow-y-auto"
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-800">
          <h2 id="share-modal-title" className="text-lg font-semibold">{t('sharing.title')}</h2>
          <button
            onClick={onClose}
            aria-label={t('common.close')}
            className="p-2 text-gray-400 hover:text-white hover:bg-gray-800 rounded-lg transition-colors"
          >
            <svg
              className="w-5 h-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6">
          {/* Preview */}
          <div className="flex justify-center overflow-x-auto">
            <ShareableProfileCard
              ref={cardRef}
              user={user}
              followerCount={followerCount}
              followingCount={followingCount}
              totalContributions={totalContributions}
              languages={languages}
              postCount={postCount}
            />
          </div>

          {/* Share buttons */}
          <ShareButtons
            onDownload={handleDownload}
            onCopyLink={handleCopyLink}
            onShareTwitter={handleShareTwitter}
            onShareLinkedIn={handleShareLinkedIn}
            generating={generating}
          />
        </div>
      </div>
    </div>
  );
}
