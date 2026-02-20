import { useTranslation } from 'react-i18next';
import type { AtCoderRatingInfo } from '../../api/atcoder';

const ATCODER_COLORS: Record<string, string> = {
  gray: '#808080', brown: '#804000', green: '#008000', cyan: '#00C0C0',
  blue: '#0000FF', yellow: '#C0C000', orange: '#FF8000', red: '#FF0000',
};

interface CompetitiveProgrammingCardProps {
  atcoderRating: AtCoderRatingInfo | null;
  atcoderUsername?: string;
  paizaRank?: string;
}

export default function CompetitiveProgrammingCard({ atcoderRating, atcoderUsername, paizaRank }: CompetitiveProgrammingCardProps) {
  const { t } = useTranslation();

  if (!atcoderRating && !paizaRank) return null;

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
      <h3 className="text-sm font-semibold text-gray-300 mb-4 flex items-center gap-2">
        <svg className="w-4 h-4 text-cyan-400" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M16.5 18.75h-9m9 0a3 3 0 0 1 3 3h-15a3 3 0 0 1 3-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497m5.007 0a7.454 7.454 0 0 1-.982-3.172M9.497 14.25a7.454 7.454 0 0 0 .981-3.172M5.25 4.236c-.982.143-1.954.317-2.916.52A6.003 6.003 0 0 0 7.73 9.728M5.25 4.236V4.5c0 2.108.966 3.99 2.48 5.228M5.25 4.236V2.721C7.456 2.41 9.71 2.25 12 2.25c2.291 0 4.545.16 6.75.47v1.516M18.75 4.236c.982.143 1.954.317 2.916.52A6.003 6.003 0 0 1 16.27 9.728M18.75 4.236V4.5c0 2.108-.966 3.99-2.48 5.228m0 0a6.023 6.023 0 0 1-2.52.587 6.023 6.023 0 0 1-2.52-.587" /></svg>
        {t('profile.competitiveProgramming')}
      </h3>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {atcoderRating && (
          <a href={`https://atcoder.jp/users/${atcoderUsername}`} target="_blank" rel="noopener noreferrer" className="flex items-center gap-4 p-4 bg-gray-800/50 rounded-lg border border-gray-700 hover:border-gray-600 transition-colors group">
            <div className="w-12 h-12 bg-gray-700 rounded-lg flex items-center justify-center text-white font-bold text-lg">A</div>
            <div>
              <div className="text-sm text-gray-400">AtCoder</div>
              <div className="flex items-center gap-2">
                <span className="text-xl font-bold" style={{ color: ATCODER_COLORS[atcoderRating.color] || '#808080' }}>
                  {atcoderRating.rating}
                </span>
                <span className="text-xs text-gray-500">({atcoderRating.rank})</span>
              </div>
            </div>
          </a>
        )}
        {paizaRank && (
          <div className="flex items-center gap-4 p-4 bg-gray-800/50 rounded-lg border border-gray-700">
            <div className="w-12 h-12 bg-emerald-700 rounded-lg flex items-center justify-center text-white font-bold text-lg">P</div>
            <div>
              <div className="text-sm text-gray-400">paiza</div>
              <div className="flex items-center gap-2">
                <span className="text-xl font-bold text-white">
                  {t('profile.paizaRankLabel', { rank: paizaRank })}
                </span>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
