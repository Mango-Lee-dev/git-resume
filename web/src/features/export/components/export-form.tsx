'use client';

import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Download, FileJson, FileText, FileSpreadsheet } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useExport } from '../hooks/use-export';
import type { ExportFormat } from '@/types';

const formats: { value: ExportFormat; label: string; icon: React.ReactNode; description: string }[] = [
  {
    value: 'csv',
    label: 'CSV',
    icon: <FileSpreadsheet className="h-5 w-5" />,
    description: 'Spreadsheet compatible format',
  },
  {
    value: 'json',
    label: 'JSON',
    icon: <FileJson className="h-5 w-5" />,
    description: 'Machine readable format',
  },
  {
    value: 'markdown',
    label: 'Markdown',
    icon: <FileText className="h-5 w-5" />,
    description: 'Documentation friendly format',
  },
];

export function ExportForm() {
  const [format, setFormat] = useState<ExportFormat>('markdown');
  const [project, setProject] = useState('');
  const [fromDate, setFromDate] = useState('');
  const [toDate, setToDate] = useState('');
  const [isExporting, setIsExporting] = useState(false);

  const { downloadExport } = useExport();

  const handleExport = async () => {
    setIsExporting(true);
    try {
      await downloadExport({
        format,
        project: project || undefined,
        from: fromDate || undefined,
        to: toDate || undefined,
      });
    } catch (error) {
      console.error('Export failed:', error);
    } finally {
      setIsExporting(false);
    }
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Export Format</CardTitle>
          <CardDescription>Choose the format for your export</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-3 sm:grid-cols-3">
            {formats.map((f) => (
              <button
                key={f.value}
                onClick={() => setFormat(f.value)}
                className={cn(
                  'flex flex-col items-center p-4 rounded-lg border-2 text-center transition-colors',
                  format === f.value
                    ? 'border-primary bg-primary/5'
                    : 'border-muted hover:border-muted-foreground/50'
                )}
              >
                {f.icon}
                <span className="font-medium mt-2">{f.label}</span>
                <span className="text-xs text-muted-foreground mt-1">
                  {f.description}
                </span>
              </button>
            ))}
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Filters (Optional)</CardTitle>
          <CardDescription>
            Narrow down the results to export
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Project</label>
            <Input
              placeholder="Filter by project name"
              value={project}
              onChange={(e) => setProject(e.target.value)}
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">From Date</label>
              <Input
                type="date"
                value={fromDate}
                onChange={(e) => setFromDate(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">To Date</label>
              <Input
                type="date"
                value={toDate}
                onChange={(e) => setToDate(e.target.value)}
              />
            </div>
          </div>
        </CardContent>
      </Card>

      <Button
        onClick={handleExport}
        disabled={isExporting}
        className="w-full"
        size="lg"
      >
        <Download className="mr-2 h-4 w-4" />
        {isExporting ? 'Exporting...' : `Export as ${format.toUpperCase()}`}
      </Button>
    </div>
  );
}
