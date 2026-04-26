export * from './categories';

export const EXPORT_FORMATS = ['csv', 'json', 'markdown'] as const;

export const TEMPLATES = [
  { name: 'default', label: 'Default', description: 'Balanced, professional, achievement-focused' },
  { name: 'startup', label: 'Startup', description: 'Dynamic, results-driven, rapid delivery emphasis' },
  { name: 'enterprise', label: 'Enterprise', description: 'Formal, process-oriented, compliance-aware' },
  { name: 'backend', label: 'Backend', description: 'Technical, systems-focused, performance/scalability' },
  { name: 'frontend', label: 'Frontend', description: 'User-centric, accessibility, design systems' },
  { name: 'devops', label: 'DevOps', description: 'Operational, metrics-driven, automation-focused' },
  { name: 'data', label: 'Data', description: 'Analytical, data-driven, data quality focus' },
] as const;
