'use client';

import { useMutation } from '@tanstack/react-query';
import { api, ENDPOINTS } from '@/lib/api';
import type { AnalyzeRequest, AnalyzeResponse } from '@/types';

export function useSubmitAnalysis() {
  return useMutation({
    mutationFn: (data: AnalyzeRequest) =>
      api.post<AnalyzeResponse>(ENDPOINTS.analyze, data),
  });
}
