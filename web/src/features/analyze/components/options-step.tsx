'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Settings } from 'lucide-react';
import { useAnalyzeStore } from '@/stores';

export function OptionsStep() {
  const { options, setOptions } = useAnalyzeStore();

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Settings className="h-5 w-5" />
          Options
        </CardTitle>
        <CardDescription>
          Configure analysis options
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-2">
          <label className="text-sm font-medium">Batch Size</label>
          <p className="text-sm text-muted-foreground">
            Number of commits to process in each API call (1-20)
          </p>
          <Input
            type="number"
            min={1}
            max={20}
            value={options.batchSize}
            onChange={(e) =>
              setOptions({ batchSize: parseInt(e.target.value) || 5 })
            }
            className="w-24"
          />
        </div>

        <div className="flex items-center justify-between">
          <div className="space-y-0.5">
            <label className="text-sm font-medium">Dry Run</label>
            <p className="text-sm text-muted-foreground">
              Preview commits without making API calls
            </p>
          </div>
          <button
            onClick={() => setOptions({ dryRun: !options.dryRun })}
            className={cn(
              'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
              options.dryRun ? 'bg-primary' : 'bg-muted'
            )}
          >
            <span
              className={cn(
                'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                options.dryRun ? 'translate-x-6' : 'translate-x-1'
              )}
            />
          </button>
        </div>
      </CardContent>
    </Card>
  );
}

function cn(...classes: (string | boolean | undefined)[]) {
  return classes.filter(Boolean).join(' ');
}
