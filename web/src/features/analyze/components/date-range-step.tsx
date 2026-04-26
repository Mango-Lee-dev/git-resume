'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { Calendar } from 'lucide-react';
import { useAnalyzeStore } from '@/stores';

const months = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December'
];

const years = Array.from({ length: 5 }, (_, i) => new Date().getFullYear() - i);

export function DateRangeStep() {
  const {
    dateMode,
    month,
    year,
    dateRange,
    setDateMode,
    setMonth,
    setYear,
    setDateRange,
  } = useAnalyzeStore();

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Calendar className="h-5 w-5" />
          Date Range
        </CardTitle>
        <CardDescription>
          Select the time period for commit analysis
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Tabs value={dateMode} onValueChange={(v) => setDateMode(v as 'month' | 'range')}>
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="month">By Month</TabsTrigger>
            <TabsTrigger value="range">Custom Range</TabsTrigger>
          </TabsList>

          <TabsContent value="month" className="space-y-4 mt-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Month</label>
                <Select
                  value={month.toString()}
                  onValueChange={(v) => v && setMonth(parseInt(v))}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {months.map((m, i) => (
                      <SelectItem key={i} value={(i + 1).toString()}>
                        {m}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Year</label>
                <Select
                  value={year.toString()}
                  onValueChange={(v) => v && setYear(parseInt(v))}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {years.map((y) => (
                      <SelectItem key={y} value={y.toString()}>
                        {y}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          </TabsContent>

          <TabsContent value="range" className="space-y-4 mt-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">From</label>
                <Input
                  type="date"
                  value={dateRange.from || ''}
                  onChange={(e) =>
                    setDateRange({ ...dateRange, from: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">To</label>
                <Input
                  type="date"
                  value={dateRange.to || ''}
                  onChange={(e) =>
                    setDateRange({ ...dateRange, to: e.target.value })
                  }
                />
              </div>
            </div>
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  );
}
