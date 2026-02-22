import { Check } from 'lucide-react';

interface Step {
  label: string;
  description?: string;
}

interface StepIndicatorProps {
  steps: Step[];
  currentStep: number;
  className?: string;
}

export default function StepIndicator({
  steps,
  currentStep,
  className = '',
}: StepIndicatorProps) {
  return (
    <div className={`flex items-start ${className}`.trim()}>
      {steps.map((step, index) => {
        const isCompleted = index < currentStep;
        const isActive = index === currentStep;

        return (
          <div key={index} className="flex items-start flex-1">
            <div className="flex flex-col items-center">
              <div
                className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium ${
                  isCompleted
                    ? 'bg-blue-600 text-white'
                    : isActive
                      ? 'bg-blue-600 text-white'
                      : 'bg-gray-700 text-gray-400'
                }`}
              >
                {isCompleted ? (
                  <Check className="w-4 h-4" />
                ) : (
                  index + 1
                )}
              </div>
              <span className={`mt-2 text-xs text-center ${isActive ? 'text-white font-medium' : 'text-gray-400'}`}>
                {step.label}
              </span>
              {step.description && (
                <span className="mt-1 text-xs text-gray-500 text-center">
                  {step.description}
                </span>
              )}
            </div>
            {index < steps.length - 1 && (
              <div
                className={`h-0.5 flex-1 mt-4 mx-2 ${
                  index < currentStep ? 'bg-blue-600' : 'bg-gray-700'
                }`}
              />
            )}
          </div>
        );
      })}
    </div>
  );
}
