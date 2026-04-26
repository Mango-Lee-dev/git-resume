'use client';

import { ENDPOINTS } from '@/lib/api';
import type { ExportQuery } from '@/types';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export function useExport() {
  const downloadExport = async (query: ExportQuery) => {
    const searchParams = new URLSearchParams();
    searchParams.set('format', query.format);
    if (query.project) searchParams.set('project', query.project);
    if (query.from) searchParams.set('from', query.from);
    if (query.to) searchParams.set('to', query.to);

    const url = `${API_BASE}${ENDPOINTS.export}?${searchParams.toString()}`;

    const response = await fetch(url);

    if (!response.ok) {
      throw new Error('Export failed');
    }

    const blob = await response.blob();
    const contentDisposition = response.headers.get('Content-Disposition');
    let filename = `export.${query.format}`;

    if (contentDisposition) {
      const match = contentDisposition.match(/filename="?(.+)"?/);
      if (match) filename = match[1];
    }

    const downloadUrl = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(downloadUrl);
  };

  return { downloadExport };
}
