'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import {
  Lock,
  Server,
  Database,
  Layout,
  TestTube,
  Zap,
  Shield,
  Cloud,
  RefreshCw,
  FileText,
  MoreHorizontal,
} from 'lucide-react';
import type { ClusterStat } from '@/types';

interface ClusterBreakdownProps {
  clusters?: Record<string, ClusterStat>;
}

const CLUSTER_CONFIG: Record<string, { icon: React.ReactNode; label: string; color: string }> = {
  auth: { icon: <Lock className="h-4 w-4" />, label: 'Authentication', color: 'bg-amber-500' },
  api: { icon: <Server className="h-4 w-4" />, label: 'API Development', color: 'bg-blue-500' },
  database: { icon: <Database className="h-4 w-4" />, label: 'Database', color: 'bg-green-500' },
  ui: { icon: <Layout className="h-4 w-4" />, label: 'UI/Frontend', color: 'bg-purple-500' },
  testing: { icon: <TestTube className="h-4 w-4" />, label: 'Testing', color: 'bg-cyan-500' },
  performance: { icon: <Zap className="h-4 w-4" />, label: 'Performance', color: 'bg-yellow-500' },
  security: { icon: <Shield className="h-4 w-4" />, label: 'Security', color: 'bg-red-500' },
  infra: { icon: <Cloud className="h-4 w-4" />, label: 'Infrastructure', color: 'bg-indigo-500' },
  refactor: { icon: <RefreshCw className="h-4 w-4" />, label: 'Refactoring', color: 'bg-orange-500' },
  docs: { icon: <FileText className="h-4 w-4" />, label: 'Documentation', color: 'bg-gray-500' },
  other: { icon: <MoreHorizontal className="h-4 w-4" />, label: 'Other', color: 'bg-slate-500' },
};

export function ClusterBreakdown({ clusters }: ClusterBreakdownProps) {
  if (!clusters || Object.keys(clusters).length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Work Pattern Clusters</CardTitle>
          <CardDescription>Similar work across projects</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">No clusters yet. Run an analysis to see patterns.</p>
        </CardContent>
      </Card>
    );
  }

  const total = Object.values(clusters).reduce((sum, stat) => sum + stat.count, 0);

  // Sort by number of projects (cross-project patterns first), then by count
  const sortedClusters = Object.entries(clusters)
    .filter(([, stat]) => stat.count > 0)
    .sort(([, a], [, b]) => {
      // Prioritize cross-project patterns
      if (b.projects.length !== a.projects.length) {
        return b.projects.length - a.projects.length;
      }
      return b.count - a.count;
    });

  return (
    <Card>
      <CardHeader>
        <CardTitle>Work Pattern Clusters</CardTitle>
        <CardDescription>Similar work across projects</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {sortedClusters.map(([cluster, stat]) => {
          const config = CLUSTER_CONFIG[cluster] || CLUSTER_CONFIG.other;
          const percentage = total > 0 ? (stat.count / total) * 100 : 0;
          const isCrossProject = stat.projects.length > 1;

          return (
            <div key={cluster} className="space-y-2">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className={`p-1 rounded ${config.color} text-white`}>
                    {config.icon}
                  </span>
                  <span className="font-medium text-sm">{config.label}</span>
                  {isCrossProject && (
                    <Badge variant="secondary" className="text-xs">
                      {stat.projects.length} projects
                    </Badge>
                  )}
                </div>
                <span className="text-sm text-muted-foreground">
                  {stat.count} ({percentage.toFixed(0)}%)
                </span>
              </div>
              <Progress value={percentage} className="h-2" />
              {isCrossProject && (
                <div className="flex flex-wrap gap-1">
                  {stat.projects.slice(0, 3).map((project) => (
                    <Badge key={project} variant="outline" className="text-xs">
                      {project}
                    </Badge>
                  ))}
                  {stat.projects.length > 3 && (
                    <Badge variant="outline" className="text-xs">
                      +{stat.projects.length - 3} more
                    </Badge>
                  )}
                </div>
              )}
            </div>
          );
        })}
      </CardContent>
    </Card>
  );
}
