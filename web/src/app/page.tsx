'use client';

import { PageHeader } from '@/components/common';
import {
  StatsCards,
  CategoryBreakdown,
  ClusterBreakdown,
  QuickActions,
  useDashboardStats,
} from '@/features/dashboard';

export default function DashboardPage() {
  const { data: stats, isLoading } = useDashboardStats();

  return (
    <div className="space-y-6">
      <PageHeader
        title="Dashboard"
        description="Overview of your Git commit analysis"
      />

      <StatsCards stats={stats} isLoading={isLoading} />

      <div className="grid gap-6 md:grid-cols-2">
        <CategoryBreakdown stats={stats} />
        <ClusterBreakdown clusters={stats?.cluster_breakdown} />
      </div>

      <QuickActions />
    </div>
  );
}
