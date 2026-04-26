'use client';

import { use } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { PageHeader } from '@/components/common';
import { JobProgress, useJobPolling } from '@/features/jobs';
import { Skeleton } from '@/components/ui/skeleton';
import { Card, CardContent } from '@/components/ui/card';

interface Props {
  params: Promise<{ jobId: string }>;
}

export default function ProgressPage({ params }: Props) {
  const { jobId } = use(params);
  const { data: job, isLoading, error } = useJobPolling(jobId);

  const isComplete = job?.status === 'completed';
  const isFailed = job?.status === 'failed' || job?.status === 'cancelled';

  return (
    <div className="max-w-2xl mx-auto">
      <PageHeader
        title="Analysis in Progress"
        description="Your commits are being analyzed"
      />

      {isLoading && (
        <Card>
          <CardContent className="py-6 space-y-4">
            <Skeleton className="h-4 w-48" />
            <Skeleton className="h-3 w-full" />
            <Skeleton className="h-3 w-32" />
          </CardContent>
        </Card>
      )}

      {error && (
        <Card>
          <CardContent className="py-6 text-center text-destructive">
            Failed to load job status. Please try refreshing.
          </CardContent>
        </Card>
      )}

      {job && <JobProgress job={job} />}

      <div className="flex justify-center gap-4 mt-6">
        {isComplete && (
          <Button asChild>
            <Link href="/results">View Results</Link>
          </Button>
        )}

        {isFailed && (
          <Button asChild variant="outline">
            <Link href="/analyze">Try Again</Link>
          </Button>
        )}

        {!isComplete && !isFailed && (
          <Button variant="outline" asChild>
            <Link href="/">Back to Dashboard</Link>
          </Button>
        )}
      </div>
    </div>
  );
}
