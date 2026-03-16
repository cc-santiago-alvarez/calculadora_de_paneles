import { CheckIcon } from '@heroicons/react/24/solid';

interface Step {
  label: string;
}

interface StepIndicatorProps {
  steps: Step[];
  currentStep: number;
}

export default function StepIndicator({ steps, currentStep }: StepIndicatorProps) {
  return (
    <div className="flex items-center gap-1">
      {steps.map((step, i) => {
        const completed = i < currentStep;
        const active = i === currentStep;

        return (
          <div key={i} className="flex items-center gap-1 flex-1">
            <div className="flex items-center gap-2 min-w-0">
              <div
                className={`w-7 h-7 rounded-full flex items-center justify-center text-xs font-medium flex-shrink-0 transition-all duration-fast ${
                  active ? 'scale-110 ' : ''
                }${
                  completed
                    ? 'bg-brand text-white'
                    : active
                      ? 'bg-brand-soft text-brand border-2 border-brand'
                      : 'bg-inset text-fg-muted border border-[var(--color-border-default)]'
                }`}
              >
                {completed ? <CheckIcon className="w-4 h-4" /> : i + 1}
              </div>
              <span
                className={`text-xs truncate hidden lg:block ${
                  active ? 'text-fg-primary font-medium' : completed ? 'text-fg-secondary' : 'text-fg-muted'
                }`}
              >
                {step.label}
              </span>
            </div>
            {i < steps.length - 1 && (
              <div className="flex-1 h-px mx-1 bg-[var(--color-border-default)] relative overflow-hidden">
                <div
                  className={`absolute inset-0 bg-brand transition-transform duration-300 ease-decel origin-left ${
                    completed ? 'scale-x-100' : 'scale-x-0'
                  }`}
                />
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
