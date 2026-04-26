'use client';

import { useQuery } from '@tanstack/react-query';
import { api, ENDPOINTS } from '@/lib/api';
import type { Job } from '@/types';

export function useJobPolling(jobId: string) {
  return useQuery({
    queryKey: ['job', jobId],
    queryFn: () => api.get<Job>(ENDPOINTS.job(jobId)),
    refetchInterval: (query) => {
      const data = query.state.data;
      if (
        data?.status === 'completed' ||
        data?.status === 'failed' ||
        data?.status === 'cancelled'
      ) {
        return false;
      }
      return 2000; // Poll every 2 seconds
    },
  });
}
