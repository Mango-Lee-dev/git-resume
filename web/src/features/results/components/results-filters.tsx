'use client';

import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { CATEGORIES } from '@/lib/constants';

interface ResultsFiltersProps {
  project: string;
  category: string;
  onProjectChange: (value: string) => void;
  onCategoryChange: (value: string) => void;
}

export function ResultsFilters({
  project,
  category,
  onProjectChange,
  onCategoryChange,
}: ResultsFiltersProps) {
  return (
    <div className="flex flex-col sm:flex-row gap-4 mb-6">
      <Input
        placeholder="Filter by project..."
        value={project}
        onChange={(e) => onProjectChange(e.target.value)}
        className="sm:w-64"
      />
      <Select value={category} onValueChange={(v) => v && onCategoryChange(v)}>
        <SelectTrigger className="sm:w-48">
          <SelectValue placeholder="All Categories" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Categories</SelectItem>
          {CATEGORIES.map((cat) => (
            <SelectItem key={cat} value={cat}>
              {cat}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
