'use client';

import { Card, CardContent } from '@/components/ui/card';
import { CategoryBadge } from './category-badge';
import { format } from 'date-fns';
import type { AnalysisResult } from '@/types';

interface ResultCardProps {
  result: AnalysisResult;
}

export function ResultCard({ result }: ResultCardProps) {
  return (
    <Card>
      <CardContent className="pt-4">
        <div className="flex items-start justify-between gap-4">
          <div className="flex-1 space-y-2">
            <p className="text-sm leading-relaxed">{result.impact_summary}</p>
            <div className="flex items-center gap-3 text-xs text-muted-foreground">
              <span>{format(new Date(result.date), 'MMM d, yyyy')}</span>
              <span>{result.project}</span>
              <code className="bg-muted px-1.5 py-0.5 rounded">
                {result.commit_hash.slice(0, 7)}
              </code>
            </div>
          </div>
          <CategoryBadge category={result.category} />
        </div>
      </CardContent>
    </Card>
  );
}
