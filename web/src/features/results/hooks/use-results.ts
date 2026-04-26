'use client';

import { useQuery } from '@tanstack/react-query';
import { api, ENDPOINTS } from '@/lib/api';
import type { AnalysisResult, PaginatedResponse, ResultsQuery } from '@/types';

export function useResults(query: ResultsQuery) {
  const searchParams = new URLSearchParams();

  if (query.page) searchParams.set('page', query.page.toString());
  if (query.page_size) searchParams.set('page_size', query.page_size.toString());
  if (query.project) searchParams.set('project', query.project);
  if (query.category) searchParams.set('category', query.category);
  if (query.from) searchParams.set('from', query.from);
  if (query.to) searchParams.set('to', query.to);

  const queryString = searchParams.toString();
  const endpoint = queryString
    ? `${ENDPOINTS.results}?${queryString}`
    : ENDPOINTS.results;

  return useQuery({
    queryKey: ['results', query],
    queryFn: () => api.get<PaginatedResponse<AnalysisResult>>(endpoint),
  });
}
