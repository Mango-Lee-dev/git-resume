'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { CATEGORY_ICONS } from '@/lib/constants';
import type { Category } from '@/types';
import type { DashboardStats } from '@/types';

interface CategoryBreakdownProps {
  stats?: DashboardStats;
}

export function CategoryBreakdown({ stats }: CategoryBreakdownProps) {
  if (!stats?.category_breakdown) {
    return null;
  }

  const breakdown = stats.category_breakdown;
  const total = Object.values(breakdown).reduce((sum, count) => sum + count, 0);

  const categories = Object.entries(breakdown)
    .filter(([, count]) => count > 0)
    .sort(([, a], [, b]) => b - a);

  if (categories.length === 0) {
    return null;
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Category Breakdown</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {categories.map(([category, count]) => {
          const percentage = total > 0 ? (count / total) * 100 : 0;
          const icon = CATEGORY_ICONS[category as Category] || '';

          return (
            <div key={category} className="space-y-1">
              <div className="flex items-center justify-between text-sm">
                <span>
                  {icon} {category}
                </span>
                <span className="text-muted-foreground">
                  {count} ({percentage.toFixed(0)}%)
                </span>
              </div>
              <Progress value={percentage} className="h-2" />
            </div>
          );
        })}
      </CardContent>
    </Card>
  );
}
