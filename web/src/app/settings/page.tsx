'use client';

import { PageHeader } from '@/components/common';
import { ApiKeyForm, SessionInfo } from '@/features/settings';
import { useSessionStore } from '@/stores';

export default function SettingsPage() {
  const { session } = useSessionStore();

  return (
    <div className="max-w-2xl mx-auto">
      <PageHeader
        title="Settings"
        description="Configure your API key and preferences"
      />

      {session ? <SessionInfo /> : <ApiKeyForm />}
    </div>
  );
}
