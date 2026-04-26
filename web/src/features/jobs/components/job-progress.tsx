'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { Badge } from '@/components/ui/badge';
import { Loader2, CheckCircle2, XCircle, Clock } from 'lucide-react';
import type { Job } from '@/types';

interface JobProgressProps {
  job: Job;
}

const phaseLabels: Record<string, string> = {
  scanning: 'Scanning repositories...',
  filtering: 'Filtering commits...',
  analyzing: 'Analyzing with Claude AI...',
  saving: 'Saving results...',
  complete: 'Analysis complete!',
  error: 'An error occurred',
};

const statusConfig: Record<string, { icon: React.ReactNode; color: string }> = {
  pending: { icon: <Clock className="h-4 w-4" />, color: 'secondary' },
  running: { icon: <Loader2 className="h-4 w-4 animate-spin" />, color: 'default' },
  completed: { icon: <CheckCircle2 className="h-4 w-4" />, color: 'default' },
  failed: { icon: <XCircle className="h-4 w-4" />, color: 'destructive' },
  cancelled: { icon: <XCircle className="h-4 w-4" />, color: 'secondary' },
};

export function JobProgress({ job }: JobProgressProps) {
  const config = statusConfig[job.status] || statusConfig.pending;
  const phaseLabel = phaseLabels[job.phase] || job.message || 'Processing...';

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Analysis Progress</CardTitle>
          <Badge variant={config.color as 'default' | 'secondary' | 'destructive'}>
            {config.icon}
            <span className="ml-1 capitalize">{job.status}</span>
          </Badge>
        </div>
        <CardDescription>{phaseLabel}</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <Progress value={job.progress} className="h-3" />

        <div className="flex justify-between text-sm text-muted-foreground">
          <span>{job.progress}% complete</span>
          {job.result_count !== undefined && job.result_count > 0 && (
            <span>{job.result_count} results generated</span>
          )}
        </div>

        {job.error && (
          <div className="p-3 bg-destructive/10 text-destructive rounded-md text-sm">
            {job.error}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
