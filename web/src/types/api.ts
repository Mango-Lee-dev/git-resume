// API response types

export interface ApiResponse<T> {
  data: T;
  message?: string;
}

export interface ApiError {
  error: string;
  code?: string;
  request_id?: string;
}

export interface PaginatedResponse<T> {
  results: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface AnalyzeResponse {
  job_id: string;
  message: string;
}

export interface ResultsQuery {
  page?: number;
  page_size?: number;
  project?: string;
  category?: string;
  from?: string;
  to?: string;
}

export interface ExportQuery {
  format: 'csv' | 'json' | 'markdown';
  project?: string;
  from?: string;
  to?: string;
}
