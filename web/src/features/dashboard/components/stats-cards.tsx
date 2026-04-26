'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { FileText, GitCommit, DollarSign, FolderGit2 } from 'lucide-react';
import type { DashboardStats } from '@/types';

interface StatsCardsProps {
  stats?: DashboardStats;
  isLoading: boolean;
}

export function StatsCards({ stats, isLoading }: StatsCardsProps) {
  const cards = [
    {
      title: 'Total Results',
      value: stats?.total_results ?? 0,
      icon: FileText,
      description: 'Generated bullet points',
    },
    {
      title: 'Commits Processed',
      value: stats?.total_commits ?? 0,
      icon: GitCommit,
      description: 'From all repositories',
    },
    {
      title: 'Projects',
      value: stats?.project_breakdown
        ? Object.keys(stats.project_breakdown).length
        : 0,
      icon: FolderGit2,
      description: 'Repositories analyzed',
    },
    {
      title: 'API Cost',
      value: stats?.tokens_used?.total_cost
        ? `$${stats.tokens_used.total_cost.toFixed(2)}`
        : '$0.00',
      icon: DollarSign,
      description: 'Total Claude API usage',
    },
  ];

  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {[...Array(4)].map((_, i) => (
          <Card key={i}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <Skeleton className="h-4 w-24" />
              <Skeleton className="h-4 w-4" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-8 w-16 mb-1" />
              <Skeleton className="h-3 w-32" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {cards.map((card) => (
        <Card key={card.title}>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">{card.title}</CardTitle>
            <card.icon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{card.value}</div>
            <p className="text-xs text-muted-foreground">{card.description}</p>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
