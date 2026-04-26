'use client';

import { useState } from 'react';
import { PageHeader, EmptyState } from '@/components/common';
import { Skeleton } from '@/components/ui/skeleton';
import { Card, CardContent } from '@/components/ui/card';
import { FileText } from 'lucide-react';
import {
  ResultCard,
  ResultsFilters,
  ResultsPagination,
  useResults,
} from '@/features/results';

export default function ResultsPage() {
  const [page, setPage] = useState(1);
  const [project, setProject] = useState('');
  const [category, setCategory] = useState('all');

  const { data, isLoading } = useResults({
    page,
    page_size: 20,
    project: project || undefined,
    category: category === 'all' ? undefined : category,
  });

  const handleProjectChange = (value: string) => {
    setProject(value);
    setPage(1);
  };

  const handleCategoryChange = (value: string) => {
    setCategory(value);
    setPage(1);
  };

  return (
    <div>
      <PageHeader
        title="Results"
        description="View your generated resume bullet points"
      />

      <ResultsFilters
        project={project}
        category={category}
        onProjectChange={handleProjectChange}
        onCategoryChange={handleCategoryChange}
      />

      {isLoading && (
        <div className="space-y-4">
          {[...Array(5)].map((_, i) => (
            <Card key={i}>
              <CardContent className="py-4 space-y-2">
                <Skeleton className="h-4 w-full" />
                <Skeleton className="h-4 w-3/4" />
                <Skeleton className="h-3 w-48" />
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {!isLoading && (!data?.results || data.results.length === 0) && (
        <EmptyState
          icon={FileText}
          title="No results yet"
          description="Start an analysis to generate resume bullet points from your Git commits."
        />
      )}

      {data?.results && data.results.length > 0 && (
        <>
          <div className="space-y-4">
            {data.results.map((result) => (
              <ResultCard key={result.id} result={result} />
            ))}
          </div>

          <ResultsPagination
            page={page}
            totalPages={data.total_pages}
            onPageChange={setPage}
          />
        </>
      )}
    </div>
  );
}
