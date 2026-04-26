import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface DateRange {
  from: string | null;
  to: string | null;
}

interface AnalyzeOptions {
  batchSize: number;
  dryRun: boolean;
}

interface AnalyzeState {
  // Wizard step
  step: number;

  // Form data
  repos: string[];
  dateMode: 'month' | 'range';
  month: number;
  year: number;
  dateRange: DateRange;
  template: string;
  options: AnalyzeOptions;

  // Current job
  currentJobId: string | null;

  // Actions
  setStep: (step: number) => void;
  nextStep: () => void;
  prevStep: () => void;
  setRepos: (repos: string[]) => void;
  addRepo: (repo: string) => void;
  removeRepo: (index: number) => void;
  setDateMode: (mode: 'month' | 'range') => void;
  setMonth: (month: number) => void;
  setYear: (year: number) => void;
  setDateRange: (range: DateRange) => void;
  setTemplate: (template: string) => void;
  setOptions: (options: Partial<AnalyzeOptions>) => void;
  setCurrentJobId: (jobId: string | null) => void;
  reset: () => void;
}

const now = new Date();
const initialState = {
  step: 0,
  repos: [],
  dateMode: 'month' as const,
  month: now.getMonth() + 1,
  year: now.getFullYear(),
  dateRange: { from: null, to: null },
  template: 'default',
  options: { batchSize: 5, dryRun: false },
  currentJobId: null,
};

export const useAnalyzeStore = create<AnalyzeState>()(
  persist(
    (set) => ({
      ...initialState,

      setStep: (step) => set({ step }),
      nextStep: () => set((state) => ({ step: state.step + 1 })),
      prevStep: () => set((state) => ({ step: Math.max(0, state.step - 1) })),

      setRepos: (repos) => set({ repos }),
      addRepo: (repo) =>
        set((state) => ({
          repos: state.repos.includes(repo) ? state.repos : [...state.repos, repo],
        })),
      removeRepo: (index) =>
        set((state) => ({
          repos: state.repos.filter((_, i) => i !== index),
        })),

      setDateMode: (dateMode) => set({ dateMode }),
      setMonth: (month) => set({ month }),
      setYear: (year) => set({ year }),
      setDateRange: (dateRange) => set({ dateRange }),

      setTemplate: (template) => set({ template }),
      setOptions: (options) =>
        set((state) => ({
          options: { ...state.options, ...options },
        })),

      setCurrentJobId: (currentJobId) => set({ currentJobId }),
      reset: () => set(initialState),
    }),
    {
      name: 'analyze-wizard',
      partialize: (state) => ({
        repos: state.repos,
        dateMode: state.dateMode,
        month: state.month,
        year: state.year,
        dateRange: state.dateRange,
        template: state.template,
        options: state.options,
      }),
    }
  )
);
