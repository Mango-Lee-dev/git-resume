// API endpoint paths

export const ENDPOINTS = {
  // Health
  health: '/health',
  ready: '/ready',

  // Analysis
  analyze: '/api/analyze',

  // Jobs
  jobs: '/api/jobs',
  job: (id: string) => `/api/jobs/${id}`,

  // Results
  results: '/api/results',
  result: (id: number) => `/api/results/${id}`,

  // Export
  export: '/api/export',

  // Templates
  templates: '/api/templates',

  // Stats
  stats: '/api/stats',
} as const;
