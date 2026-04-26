'use client';

import { useMutation } from '@tanstack/react-query';
import { api, ENDPOINTS } from '@/lib/api';

interface CreateSessionRequest {
  api_key: string;
}

interface CreateSessionResponse {
  session_id: string;
  created_at: string;
  expires_at: string;
  message: string;
}

export function useCreateSession() {
  return useMutation({
    mutationFn: (data: CreateSessionRequest) =>
      api.post<CreateSessionResponse>(ENDPOINTS.sessions, data),
  });
}
