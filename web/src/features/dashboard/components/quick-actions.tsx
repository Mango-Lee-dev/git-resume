'use client';

import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { PlayCircle, FileText, Download } from 'lucide-react';

export function QuickActions() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Quick Actions</CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        <Button asChild className="w-full justify-start" variant="outline">
          <Link href="/analyze">
            <PlayCircle className="mr-2 h-4 w-4" />
            Start New Analysis
          </Link>
        </Button>
        <Button asChild className="w-full justify-start" variant="outline">
          <Link href="/results">
            <FileText className="mr-2 h-4 w-4" />
            View All Results
          </Link>
        </Button>
        <Button asChild className="w-full justify-start" variant="outline">
          <Link href="/export">
            <Download className="mr-2 h-4 w-4" />
            Export Results
          </Link>
        </Button>
      </CardContent>
    </Card>
  );
}
