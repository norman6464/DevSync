import { useCallback } from 'react';
import { Navigate } from 'react-router-dom';
import { User, Code, Link, CheckCircle } from 'lucide-react';
import { sectionContainerClass } from '../constants/styles';
import { useOnboarding } from '../hooks/useOnboarding';
import OnboardingProfileStep from '../components/onboarding/OnboardingProfileStep';
import OnboardingSkillsStep from '../components/onboarding/OnboardingSkillsStep';
import OnboardingIntegrationsStep from '../components/onboarding/OnboardingIntegrationsStep';
import OnboardingCompleteStep from '../components/onboarding/OnboardingCompleteStep';

const STEPS = [
  { id: 1, icon: User },
  { id: 2, icon: Code },
  { id: 3, icon: Link },
  { id: 4, icon: CheckCircle },
];

export default function OnboardingPage() {
  const {
    user,
    step,
    setStep,
    saving,
    name,
    setName,
    bio,
    setBio,
    handleSaveProfile,
    selectedLanguages,
    selectedFrameworks,
    toggleLanguage,
    toggleFramework,
    handleSaveSkills,
    zennUsername,
    setZennUsername,
    qiitaUsername,
    setQiitaUsername,
    atcoderUsername,
    setAtcoderUsername,
    connectingZenn,
    connectingQiita,
    connectingAtcoder,
    paizaRank,
    setPaizaRank,
    handleConnectGitHub,
    handleConnectZenn,
    handleConnectQiita,
    handleConnectAtCoder,
    handleSavePaizaRank,
    handleComplete,
  } = useOnboarding();

  const handleGoToStep1 = useCallback(() => setStep(1), [setStep]);
  const handleGoToStep2 = useCallback(() => setStep(2), [setStep]);
  const handleGoToStep3 = useCallback(() => setStep(3), [setStep]);
  const handleGoToStep4 = useCallback(() => setStep(4), [setStep]);

  if (!user) return null;
  if (user.onboarding_completed) return <Navigate to="/" replace />;

  return (
    <div className="min-h-screen bg-gray-950 flex flex-col items-center justify-center px-4 py-12">
      <div className="w-full max-w-xl">
        {/* Logo */}
        <div className="text-center mb-8">
          <svg className="w-12 h-12 mx-auto mb-4 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M16 18l2-2-2-2" />
            <path d="M8 6L6 8l2 2" />
            <path d="M14.5 4l-5 16" />
          </svg>
        </div>

        {/* Step Indicator */}
        <div className="flex items-center justify-center gap-2 mb-8">
          {STEPS.map((s, i) => {
            const StepIcon = s.icon;
            return (
              <div key={s.id} className="flex items-center">
                <div
                  className={`w-10 h-10 rounded-full flex items-center justify-center border-2 transition-colors ${
                    step === s.id
                      ? 'border-blue-500 bg-blue-500/20 text-blue-400'
                      : step > s.id
                      ? 'border-green-500 bg-green-500/20 text-green-400'
                      : 'border-gray-700 bg-gray-800 text-gray-500'
                  }`}
                >
                  <StepIcon className="w-5 h-5" />
                </div>
                {i < STEPS.length - 1 && (
                  <div
                    className={`w-12 h-0.5 mx-1 transition-colors ${
                      step > s.id ? 'bg-green-500' : 'bg-gray-700'
                    }`}
                  />
                )}
              </div>
            );
          })}
        </div>

        {/* Step Content */}
        <div className={sectionContainerClass}>
          {step === 1 && (
            <OnboardingProfileStep
              name={name}
              setName={setName}
              bio={bio}
              setBio={setBio}
              saving={saving}
              onSave={handleSaveProfile}
              onSkip={handleGoToStep2}
            />
          )}
          {step === 2 && (
            <OnboardingSkillsStep
              selectedLanguages={selectedLanguages}
              selectedFrameworks={selectedFrameworks}
              toggleLanguage={toggleLanguage}
              toggleFramework={toggleFramework}
              saving={saving}
              onSave={handleSaveSkills}
              onBack={handleGoToStep1}
              onSkip={handleGoToStep3}
            />
          )}
          {step === 3 && (
            <OnboardingIntegrationsStep
              user={user}
              zennUsername={zennUsername}
              setZennUsername={setZennUsername}
              qiitaUsername={qiitaUsername}
              setQiitaUsername={setQiitaUsername}
              atcoderUsername={atcoderUsername}
              setAtcoderUsername={setAtcoderUsername}
              paizaRank={paizaRank}
              setPaizaRank={setPaizaRank}
              connectingZenn={connectingZenn}
              connectingQiita={connectingQiita}
              connectingAtcoder={connectingAtcoder}
              onConnectGitHub={handleConnectGitHub}
              onConnectZenn={handleConnectZenn}
              onConnectQiita={handleConnectQiita}
              onConnectAtCoder={handleConnectAtCoder}
              onSavePaizaRank={handleSavePaizaRank}
              onBack={handleGoToStep2}
              onNext={handleGoToStep4}
            />
          )}
          {step === 4 && (
            <OnboardingCompleteStep
              saving={saving}
              onComplete={handleComplete}
            />
          )}
        </div>
      </div>
    </div>
  );
}
