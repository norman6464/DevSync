import { useTranslation } from 'react-i18next';
import {
  Sprout, Monitor, Swords, Crown, Trophy, Flame, Zap,
  PenLine, FileText, Heart, Star, Handshake, Megaphone, StarIcon,
  Medal, MessageCircleQuestion, HelpCircle, Award, GraduationCap,
  type LucideIcon,
} from 'lucide-react';
import type { BadgeResult } from '../../types/badge';

interface BadgeDisplayProps {
  badges: BadgeResult[];
}

/** Badge ID → icon, color, bgColor mapping (UI responsibility) */
const badgeMeta: Record<string, { icon: LucideIcon; color: string; bgColor: string }> = {
  'first-commit':    { icon: Sprout,     color: 'text-green-400',  bgColor: 'bg-green-500/10 border-green-500/30' },
  'contributor':     { icon: Monitor,    color: 'text-blue-400',   bgColor: 'bg-blue-500/10 border-blue-500/30' },
  'code-warrior':    { icon: Swords,     color: 'text-purple-400', bgColor: 'bg-purple-500/10 border-purple-500/30' },
  'commit-master':   { icon: Crown,      color: 'text-yellow-400', bgColor: 'bg-yellow-500/10 border-yellow-500/30' },
  'legend':          { icon: Trophy,     color: 'text-orange-400', bgColor: 'bg-orange-500/10 border-orange-500/30' },
  'week-streak':     { icon: Flame,      color: 'text-orange-400', bgColor: 'bg-orange-500/10 border-orange-500/30' },
  'month-streak':    { icon: Zap,        color: 'text-yellow-400', bgColor: 'bg-yellow-500/10 border-yellow-500/30' },
  'first-post':      { icon: PenLine,    color: 'text-cyan-400',   bgColor: 'bg-cyan-500/10 border-cyan-500/30' },
  'blogger':         { icon: FileText,   color: 'text-indigo-400', bgColor: 'bg-indigo-500/10 border-indigo-500/30' },
  'liked':           { icon: Heart,      color: 'text-pink-400',   bgColor: 'bg-pink-500/10 border-pink-500/30' },
  'popular':         { icon: Star,       color: 'text-rose-400',   bgColor: 'bg-rose-500/10 border-rose-500/30' },
  'friendly':        { icon: Handshake,  color: 'text-teal-400',   bgColor: 'bg-teal-500/10 border-teal-500/30' },
  'influencer':      { icon: Megaphone,  color: 'text-violet-400', bgColor: 'bg-violet-500/10 border-violet-500/30' },
  'star':            { icon: StarIcon,   color: 'text-amber-400',  bgColor: 'bg-amber-500/10 border-amber-500/30' },
  'qa-first-answer': { icon: MessageCircleQuestion, color: 'text-blue-400', bgColor: 'bg-blue-500/10 border-blue-500/30' },
  'qa-helper':       { icon: HelpCircle, color: 'text-cyan-400',   bgColor: 'bg-cyan-500/10 border-cyan-500/30' },
  'goal-achiever':   { icon: Award,      color: 'text-green-400',  bgColor: 'bg-green-500/10 border-green-500/30' },
  'goal-master':     { icon: GraduationCap, color: 'text-purple-400', bgColor: 'bg-purple-500/10 border-purple-500/30' },
};

const fallbackMeta = { icon: Medal, color: 'text-gray-400', bgColor: 'bg-gray-500/10 border-gray-500/30' };

export default function BadgeDisplay({ badges }: BadgeDisplayProps) {
  const { t } = useTranslation();

  const earnedBadges = badges.filter((b) => b.earned);
  const lockedBadges = badges.filter((b) => !b.earned);

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-800 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <svg
            className="w-4 h-4 text-yellow-400"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M16.5 18.75h-9m9 0a3 3 0 0 1 3 3h-15a3 3 0 0 1 3-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497m5.007 0a7.454 7.454 0 0 1-.982-3.172M9.497 14.25a7.454 7.454 0 0 0 .981-3.172M5.25 4.236c-.982.143-1.954.317-2.916.52A6.003 6.003 0 0 0 7.73 9.728M5.25 4.236V4.5c0 2.108.966 3.99 2.48 5.228M5.25 4.236V2.721C7.456 2.41 9.71 2.25 12 2.25c2.291 0 4.545.16 6.75.47v1.516M7.73 9.728a6.726 6.726 0 0 0 2.748 1.35m8.272-6.842V4.5c0 2.108-.966 3.99-2.48 5.228m2.48-5.492a46.32 46.32 0 0 1 2.916.52 6.003 6.003 0 0 1-5.395 4.972m0 0a6.726 6.726 0 0 1-2.749 1.35m0 0a6.772 6.772 0 0 1-2.748 0"
            />
          </svg>
          <h2 className="text-sm font-semibold">{t('profile.achievements')}</h2>
        </div>
        <span className="text-xs text-gray-500">
          {earnedBadges.length}/{badges.length} {t('profile.unlocked')}
        </span>
      </div>

      <div className="p-6">
        {/* Earned Badges */}
        {earnedBadges.length > 0 && (
          <div className="mb-6">
            <h3 className="text-xs text-gray-400 uppercase tracking-wide mb-3">{t('profile.earned')}</h3>
            <div className="flex flex-wrap gap-2">
              {earnedBadges.map((badge) => {
                const meta = badgeMeta[badge.id] || fallbackMeta;
                const Icon = meta.icon;
                return (
                  <div
                    key={badge.id}
                    className={`group relative px-3 py-2 rounded-lg border ${meta.bgColor} cursor-pointer transition-all hover:scale-105`}
                  >
                    <div className="flex items-center gap-2">
                      <Icon className={`w-5 h-5 ${meta.color}`} />
                      <span className={`text-sm font-medium ${meta.color}`}>{t(badge.name)}</span>
                    </div>
                    <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-gray-800 rounded-lg text-xs text-center opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all whitespace-nowrap z-10 shadow-lg">
                      <div className="font-medium text-white">{t(badge.name)}</div>
                      <div className="text-gray-400">{t(badge.description)}</div>
                      <div className="absolute top-full left-1/2 -translate-x-1/2 border-4 border-transparent border-t-gray-800" />
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}

        {/* Locked Badges */}
        {lockedBadges.length > 0 && (
          <div>
            <h3 className="text-xs text-gray-500 uppercase tracking-wide mb-3">{t('profile.locked')}</h3>
            <div className="flex flex-wrap gap-2">
              {lockedBadges.map((badge) => {
                const meta = badgeMeta[badge.id] || fallbackMeta;
                const Icon = meta.icon;
                return (
                  <div
                    key={badge.id}
                    className="group relative px-3 py-2 rounded-lg border border-gray-700/50 bg-gray-800/30 cursor-pointer transition-all hover:border-gray-600"
                  >
                    <div className="flex items-center gap-2 opacity-40">
                      <Icon className="w-5 h-5 text-gray-500" />
                      <span className="text-sm font-medium text-gray-500">{t(badge.name)}</span>
                    </div>
                    <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-gray-800 rounded-lg text-xs text-center opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all whitespace-nowrap z-10 shadow-lg">
                      <div className="font-medium text-white">{t(badge.name)}</div>
                      <div className="text-gray-400">{t(badge.description)}</div>
                      <div className="absolute top-full left-1/2 -translate-x-1/2 border-4 border-transparent border-t-gray-800" />
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}

        {/* Empty State */}
        {earnedBadges.length === 0 && (
          <div className="text-center text-gray-500 py-4">
            <Medal className="w-8 h-8 text-gray-500 mx-auto mb-2" />
            <p className="text-sm">{t('profile.startContributing')}</p>
          </div>
        )}
      </div>
    </div>
  );
}
