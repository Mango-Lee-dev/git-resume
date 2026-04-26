// Domain models - mirrors Go pkg/models

export type Category = 'Feature' | 'Fix' | 'Refactor' | 'Test' | 'Docs' | 'Chore';

export type ExportFormat = 'csv' | 'json' | 'markdown';

export type JobStatus = 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';

export type JobPhase = 'scanning' | 'filtering' | 'analyzing' | 'saving' | 'complete' | 'error';

export interface AnalysisResult {
  id: number;
  commit_hash: string;
  date: string;
  project: string;
  category: Category;
  impact_summary: string;
  created_at: string;
}

export interface Template {
  name: string;
  description: string;
  persona: string;
  tone_style: string;
  focus: string[];
  keywords: Record<string, string>;
  output_hints: string[];
}

export interface Job {
  id: string;
  status: JobStatus;
  progress: number;
  phase: JobPhase;
  message: string;
  created_at: string;
  started_at?: string;
  completed_at?: string;
  result_count?: number;
  error?: string;
}

export interface TokenUsage {
  input_tokens: number;
  output_tokens: number;
  total_cost: number;
}

export interface DashboardStats {
  total_results: number;
  total_commits: number;
  tokens_used: TokenUsage;
  category_breakdown: Record<Category, number>;
  project_breakdown: Record<string, number>;
  recent_activity: Array<{ date: string; count: number }>;
}

export interface AnalyzeRequest {
  repos: string[];
  from_date?: string;
  to_date?: string;
  month?: number;
  year?: number;
  template: string;
  batch_size: number;
  dry_run: boolean;
}
