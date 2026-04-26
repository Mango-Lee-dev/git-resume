'use client';

import { useQuery } from '@tanstack/react-query';
import { api, ENDPOINTS } from '@/lib/api';
import type { DashboardStats } from '@/types';

export function useDashboardStats() {
  return useQuery({
    queryKey: ['stats'],
    queryFn: () => api.get<DashboardStats>(ENDPOINTS.stats),
  });
}
