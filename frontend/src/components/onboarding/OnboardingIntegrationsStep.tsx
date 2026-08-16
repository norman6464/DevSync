import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { CheckCircle, ChevronLeft, ChevronRight } from 'lucide-react';
import { inputClass, buttonPrimaryClass } from '../../constants/styles';
import type { User } from '../../types/user';
import IntegrationUsernameCard from './IntegrationUsernameCard';

interface OnboardingIntegrationsStepProps {
  user: User;
  zennUsername: string;
  setZennUsername: (v: string) => void;
  qiitaUsername: string;
  setQiitaUsername: (v: string) => void;
  atcoderUsername: string;
  setAtcoderUsername: (v: string) => void;
  paizaRank: string;
  setPaizaRank: (v: string) => void;
  connectingZenn: boolean;
  connectingQiita: boolean;
  connectingAtcoder: boolean;
  onConnectGitHub: () => void;
  onConnectZenn: () => void;
  onConnectQiita: () => void;
  onConnectAtCoder: () => void;
  onSavePaizaRank: () => void;
  onBack: () => void;
  onNext: () => void;
}

export default function OnboardingIntegrationsStep({
  user,
  zennUsername,
  setZennUsername,
  qiitaUsername,
  setQiitaUsername,
  atcoderUsername,
  setAtcoderUsername,
  paizaRank,
  setPaizaRank,
  connectingZenn,
  connectingQiita,
  connectingAtcoder,
  onConnectGitHub,
  onConnectZenn,
  onConnectQiita,
  onConnectAtCoder,
  onSavePaizaRank,
  onBack,
  onNext,
}: OnboardingIntegrationsStepProps) {
  const { t } = useTranslation();
  const handleZennChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setZennUsername(e.target.value), [setZennUsername]);
  const handleQiitaChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setQiitaUsername(e.target.value), [setQiitaUsername]);
  const handleAtcoderChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setAtcoderUsername(e.target.value), [setAtcoderUsername]);
  const handlePaizaChange = useCallback((e: React.ChangeEvent<HTMLSelectElement>) => setPaizaRank(e.target.value), [setPaizaRank]);

  return (
    <div>
      <div className="px-6 py-5 border-b border-gray-800">
        <h2 className="text-xl font-semibold text-white">{t('onboarding.integrationsTitle')}</h2>
        <p className="text-sm text-gray-400 mt-1">{t('onboarding.integrationsDescription')}</p>
      </div>
      <div className="p-6 space-y-4">
        {/* GitHub */}
        <div className="p-4 bg-gray-800/50 rounded-lg border border-gray-700">
          <div className="flex items-center gap-3 mb-3">
            <svg className="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
            <div>
              <h3 className="text-sm font-medium text-white">GitHub</h3>
              <p className="text-xs text-gray-400">{t('onboarding.githubDescription')}</p>
            </div>
          </div>
          {user.github_connected ? (
            <div className="flex items-center gap-2 text-green-400 text-sm">
              <CheckCircle className="w-4 h-4" />
              <span>{t('settings.connected')} - @{user.github_username}</span>
            </div>
          ) : (
            <button
              onClick={onConnectGitHub}
              className="w-full py-2 bg-white hover:bg-gray-100 text-gray-900 rounded-lg font-medium text-sm transition-colors inline-flex items-center justify-center gap-2"
            >
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
              {t('settings.connect')} GitHub
            </button>
          )}
        </div>

        {/* Zenn */}
        <IntegrationUsernameCard
          icon={<div className="w-6 h-6 bg-blue-500 rounded flex items-center justify-center text-white font-bold text-xs">Z</div>}
          serviceName="Zenn"
          description={t('onboarding.zennDescription')}
          connectedUsername={user.zenn_username}
          username={zennUsername}
          onUsernameChange={handleZennChange}
          placeholder={t('settings.zennUsername')}
          connecting={connectingZenn}
          onConnect={onConnectZenn}
          buttonClassName={`${buttonPrimaryClass} text-sm`}
        />

        {/* Qiita */}
        <IntegrationUsernameCard
          icon={<div className="w-6 h-6 bg-green-500 rounded flex items-center justify-center text-white font-bold text-xs">Q</div>}
          serviceName="Qiita"
          description={t('onboarding.qiitaDescription')}
          connectedUsername={user.qiita_username}
          username={qiitaUsername}
          onUsernameChange={handleQiitaChange}
          placeholder={t('settings.qiitaUsername')}
          connecting={connectingQiita}
          onConnect={onConnectQiita}
          buttonClassName="px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg font-medium text-sm transition-colors"
        />

        {/* AtCoder */}
        <IntegrationUsernameCard
          icon={<div className="w-6 h-6 bg-gray-700 rounded flex items-center justify-center text-white font-bold text-xs">A</div>}
          serviceName="AtCoder"
          description={t('onboarding.atcoderDescription')}
          connectedUsername={user.atcoder_username}
          username={atcoderUsername}
          onUsernameChange={handleAtcoderChange}
          placeholder={t('settings.atcoderUsername')}
          connecting={connectingAtcoder}
          onConnect={onConnectAtCoder}
          buttonClassName="px-4 py-2 bg-cyan-600 hover:bg-cyan-500 text-white rounded-lg font-medium text-sm transition-colors"
        />

        {/* paiza */}
        <div className="p-4 bg-gray-800/50 rounded-lg border border-gray-700">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-6 h-6 bg-emerald-700 rounded flex items-center justify-center text-white font-bold text-xs">P</div>
            <div>
              <h3 className="text-sm font-medium text-white">paiza</h3>
              <p className="text-xs text-gray-400">{t('onboarding.paizaDescription')}</p>
            </div>
          </div>
          {user.paiza_rank ? (
            <div className="flex items-center gap-2 text-green-400 text-sm">
              <CheckCircle className="w-4 h-4" />
              <span>{t('settings.connected')} - {t('settings.paizaRankLabel')}: {user.paiza_rank}</span>
            </div>
          ) : (
            <div className="flex gap-2">
              <select
                value={paizaRank}
                onChange={handlePaizaChange}
                className={`${inputClass} flex-1`}
              >
                <option value="">{t('settings.paizaSelectRank')}</option>
                <option value="S">S</option>
                <option value="A">A</option>
                <option value="B">B</option>
                <option value="C">C</option>
                <option value="D">D</option>
                <option value="E">E</option>
              </select>
              <button
                onClick={onSavePaizaRank}
                disabled={!paizaRank}
                className="px-4 py-2 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white rounded-lg font-medium text-sm transition-colors whitespace-nowrap"
              >
                {t('common.save')}
              </button>
            </div>
          )}
        </div>
      </div>
      <div className="px-6 py-4 border-t border-gray-800 flex justify-between">
        <button
          onClick={onBack}
          className="px-4 py-2 text-gray-400 hover:text-white text-sm transition-colors inline-flex items-center gap-1"
        >
          <ChevronLeft className="w-4 h-4" />
          {t('onboarding.back')}
        </button>
        <button
          onClick={onNext}
          className="px-5 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg font-medium text-sm transition-colors inline-flex items-center gap-1"
        >
          {t('onboarding.next')}
          <ChevronRight className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
}
