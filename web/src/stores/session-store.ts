'use client';

import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface Session {
  id: string;
  createdAt: string;
  expiresAt: string;
}

interface SessionState {
  apiKey: string | null;
  session: Session | null;
  isValidating: boolean;

  setApiKey: (key: string | null) => void;
  setSession: (session: Session | null) => void;
  setIsValidating: (validating: boolean) => void;
  clearSession: () => void;
  isAuthenticated: () => boolean;
}

export const useSessionStore = create<SessionState>()(
  persist(
    (set, get) => ({
      apiKey: null,
      session: null,
      isValidating: false,

      setApiKey: (apiKey) => set({ apiKey }),
      setSession: (session) => set({ session }),
      setIsValidating: (isValidating) => set({ isValidating }),

      clearSession: () =>
        set({
          apiKey: null,
          session: null,
        }),

      isAuthenticated: () => {
        const { session } = get();
        if (!session) return false;

        // Check if session is expired
        const expiresAt = new Date(session.expiresAt);
        return expiresAt > new Date();
      },
    }),
    {
      name: 'git-resume-session',
      partialize: (state) => ({
        apiKey: state.apiKey,
        session: state.session,
      }),
    }
  )
);
