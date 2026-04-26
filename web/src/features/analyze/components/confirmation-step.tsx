'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { CheckCircle2 } from 'lucide-react';
import { useAnalyzeStore } from '@/stores';
import { TEMPLATES } from '@/lib/constants';

const months = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December'
];

export function ConfirmationStep() {
  const { repos, dateMode, month, year, dateRange, template, options } =
    useAnalyzeStore();

  const templateInfo = TEMPLATES.find((t) => t.name === template);

  const getDateDisplay = () => {
    if (dateMode === 'month') {
      return `${months[month - 1]} ${year}`;
    }
    return `${dateRange.from || 'Start'} - ${dateRange.to || 'End'}`;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <CheckCircle2 className="h-5 w-5" />
          Confirm Settings
        </CardTitle>
        <CardDescription>
          Review your analysis configuration before starting
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <h4 className="text-sm font-medium">Repositories</h4>
          <div className="flex flex-wrap gap-2">
            {repos.map((repo, i) => (
              <Badge key={i} variant="secondary" className="font-mono">
                {repo.split('/').pop()}
              </Badge>
            ))}
          </div>
        </div>

        <div className="space-y-2">
          <h4 className="text-sm font-medium">Date Range</h4>
          <p className="text-sm text-muted-foreground">{getDateDisplay()}</p>
        </div>

        <div className="space-y-2">
          <h4 className="text-sm font-medium">Template</h4>
          <p className="text-sm text-muted-foreground">
            {templateInfo?.label} - {templateInfo?.description}
          </p>
        </div>

        <div className="space-y-2">
          <h4 className="text-sm font-medium">Options</h4>
          <div className="flex gap-2">
            <Badge variant="outline">Batch Size: {options.batchSize}</Badge>
            {options.dryRun && <Badge variant="outline">Dry Run</Badge>}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
