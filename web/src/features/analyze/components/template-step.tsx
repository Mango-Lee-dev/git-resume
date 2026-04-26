'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { FileCode2 } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useAnalyzeStore } from '@/stores';
import { TEMPLATES } from '@/lib/constants';

export function TemplateStep() {
  const { template, setTemplate } = useAnalyzeStore();

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <FileCode2 className="h-5 w-5" />
          Select Template
        </CardTitle>
        <CardDescription>
          Choose a template that matches your target role or company style
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid gap-3 sm:grid-cols-2">
          {TEMPLATES.map((t) => (
            <button
              key={t.name}
              onClick={() => setTemplate(t.name)}
              className={cn(
                'flex flex-col items-start p-4 rounded-lg border-2 text-left transition-colors',
                template === t.name
                  ? 'border-primary bg-primary/5'
                  : 'border-muted hover:border-muted-foreground/50'
              )}
            >
              <span className="font-medium">{t.label}</span>
              <span className="text-sm text-muted-foreground mt-1">
                {t.description}
              </span>
            </button>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
