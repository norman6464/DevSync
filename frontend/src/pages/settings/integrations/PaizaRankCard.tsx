import { useTranslation } from 'react-i18next';
import { sectionContainerClass, labelClass, selectClass } from '../../../constants/styles';

interface Props {
  paizaRank: string;
  setPaizaRank: (v: string) => void;
  saving: boolean;
  onSave: () => void;
}

export default function PaizaRankCard({ paizaRank, setPaizaRank, saving, onSave }: Props) {
  const { t } = useTranslation();

  return (
    <div className={sectionContainerClass}>
      <div className="px-6 py-4 border-b border-gray-800">
        <h2 className="text-base font-semibold">{t('settings.paiza')}</h2>
      </div>
      <div className="p-6 space-y-4">
        <div className="text-center py-2">
          <div className="w-12 h-12 bg-emerald-700/50 rounded-lg flex items-center justify-center text-emerald-400 font-bold text-xl mx-auto mb-3">P</div>
          <p className="text-gray-400 text-sm mb-4">{t('settings.paizaDescription')}</p>
        </div>
        <div>
          <label htmlFor="integration-paiza-rank" className={labelClass}>{t('settings.paizaRankLabel')}</label>
          <select
            id="integration-paiza-rank"
            value={paizaRank}
            onChange={(e) => setPaizaRank(e.target.value)}
            className={`${selectClass} w-full`}
          >
            <option value="">{t('settings.paizaSelectRank')}</option>
            <option value="S">S {t('settings.paizaRankS')}</option>
            <option value="A">A {t('settings.paizaRankA')}</option>
            <option value="B">B {t('settings.paizaRankB')}</option>
            <option value="C">C {t('settings.paizaRankC')}</option>
            <option value="D">D {t('settings.paizaRankD')}</option>
            <option value="E">E {t('settings.paizaRankE')}</option>
          </select>
        </div>
        <button
          type="button"
          onClick={onSave}
          disabled={saving}
          className="w-full px-5 py-2.5 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white rounded-lg font-semibold text-sm transition-colors"
        >
          {saving ? t('common.loading') : t('common.save')}
        </button>
      </div>
    </div>
  );
}
