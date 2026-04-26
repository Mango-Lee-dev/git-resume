'use client';

import { PageHeader } from '@/components/common';
import { ExportForm } from '@/features/export';

export default function ExportPage() {
  return (
    <div className="max-w-2xl mx-auto">
      <PageHeader
        title="Export Results"
        description="Download your resume bullet points in various formats"
      />
      <ExportForm />
    </div>
  );
}
