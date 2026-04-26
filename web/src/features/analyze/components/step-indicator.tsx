'use client';

import { cn } from '@/lib/utils';
import { Check } from 'lucide-react';

interface Step {
  id: number;
  name: string;
}

interface StepIndicatorProps {
  steps: Step[];
  currentStep: number;
}

export function StepIndicator({ steps, currentStep }: StepIndicatorProps) {
  return (
    <nav aria-label="Progress" className="mb-8">
      <ol className="flex items-center justify-between">
        {steps.map((step, index) => (
          <li key={step.id} className="flex items-center">
            <div className="flex items-center">
              <span
                className={cn(
                  'flex h-8 w-8 items-center justify-center rounded-full text-sm font-medium',
                  currentStep > index
                    ? 'bg-primary text-primary-foreground'
                    : currentStep === index
                    ? 'border-2 border-primary text-primary'
                    : 'border-2 border-muted text-muted-foreground'
                )}
              >
                {currentStep > index ? (
                  <Check className="h-4 w-4" />
                ) : (
                  index + 1
                )}
              </span>
              <span
                className={cn(
                  'ml-2 text-sm font-medium hidden sm:block',
                  currentStep >= index
                    ? 'text-foreground'
                    : 'text-muted-foreground'
                )}
              >
                {step.name}
              </span>
            </div>
            {index < steps.length - 1 && (
              <div
                className={cn(
                  'mx-4 h-0.5 w-12 sm:w-24',
                  currentStep > index ? 'bg-primary' : 'bg-muted'
                )}
              />
            )}
          </li>
        ))}
      </ol>
    </nav>
  );
}
