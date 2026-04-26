'use client';

import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { PageHeader } from '@/components/common';
import {
  StepIndicator,
  RepoStep,
  DateRangeStep,
  TemplateStep,
  OptionsStep,
  ConfirmationStep,
  useSubmitAnalysis,
} from '@/features/analyze';
import { useAnalyzeStore } from '@/stores';
import type { AnalyzeRequest } from '@/types';

const STEPS = [
  { id: 0, name: 'Repository' },
  { id: 1, name: 'Date Range' },
  { id: 2, name: 'Template' },
  { id: 3, name: 'Options' },
  { id: 4, name: 'Confirm' },
];

export default function AnalyzePage() {
  const router = useRouter();
  const {
    step,
    repos,
    dateMode,
    month,
    year,
    dateRange,
    template,
    options,
    nextStep,
    prevStep,
    setCurrentJobId,
    reset,
  } = useAnalyzeStore();

  const { mutate: submitAnalysis, isPending } = useSubmitAnalysis();

  const canProceed = () => {
    switch (step) {
      case 0:
        return repos.length > 0;
      case 1:
        if (dateMode === 'month') return true;
        return dateRange.from && dateRange.to;
      case 2:
        return !!template;
      case 3:
        return options.batchSize >= 1 && options.batchSize <= 20;
      case 4:
        return true;
      default:
        return false;
    }
  };

  const handleSubmit = () => {
    const request: AnalyzeRequest = {
      repos,
      template,
      batch_size: options.batchSize,
      dry_run: options.dryRun,
    };

    if (dateMode === 'month') {
      request.month = month;
      request.year = year;
    } else {
      request.from_date = dateRange.from || undefined;
      request.to_date = dateRange.to || undefined;
    }

    submitAnalysis(request, {
      onSuccess: (data) => {
        setCurrentJobId(data.job_id);
        router.push(`/analyze/progress/${data.job_id}`);
      },
    });
  };

  const renderStep = () => {
    switch (step) {
      case 0:
        return <RepoStep />;
      case 1:
        return <DateRangeStep />;
      case 2:
        return <TemplateStep />;
      case 3:
        return <OptionsStep />;
      case 4:
        return <ConfirmationStep />;
      default:
        return null;
    }
  };

  return (
    <div className="max-w-2xl mx-auto">
      <PageHeader
        title="Analyze Commits"
        description="Configure and start a new Git commit analysis"
      />

      <StepIndicator steps={STEPS} currentStep={step} />

      <div className="mb-6">{renderStep()}</div>

      <div className="flex justify-between">
        <Button
          variant="outline"
          onClick={step === 0 ? reset : prevStep}
          disabled={isPending}
        >
          {step === 0 ? 'Reset' : 'Back'}
        </Button>

        {step < STEPS.length - 1 ? (
          <Button onClick={nextStep} disabled={!canProceed()}>
            Next
          </Button>
        ) : (
          <Button onClick={handleSubmit} disabled={!canProceed() || isPending}>
            {isPending ? 'Starting...' : 'Start Analysis'}
          </Button>
        )}
      </div>
    </div>
  );
}
