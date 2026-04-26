'use client';

import { useRouter } from 'next/navigation';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Key, LogOut, Clock } from 'lucide-react';
import { useSessionStore } from '@/stores';
import { format } from 'date-fns';

export function SessionInfo() {
  const router = useRouter();
  const { apiKey, session, clearSession } = useSessionStore();

  const handleLogout = () => {
    clearSession();
    router.push('/settings');
  };

  if (!session || !apiKey) {
    return null;
  }

  const maskedKey = `${apiKey.slice(0, 10)}...${apiKey.slice(-4)}`;
  const expiresAt = new Date(session.expiresAt);
  const isExpired = expiresAt < new Date();

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Key className="h-5 w-5" />
          Current Session
        </CardTitle>
        <CardDescription>Your active API session</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">API Key</span>
            <code className="text-sm bg-muted px-2 py-1 rounded">{maskedKey}</code>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Session ID</span>
            <code className="text-sm bg-muted px-2 py-1 rounded">
              {session.id.slice(0, 8)}...
            </code>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Status</span>
            <Badge variant={isExpired ? 'destructive' : 'default'}>
              {isExpired ? 'Expired' : 'Active'}
            </Badge>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Expires</span>
            <span className="text-sm flex items-center gap-1">
              <Clock className="h-3 w-3" />
              {format(expiresAt, 'MMM d, yyyy HH:mm')}
            </span>
          </div>
        </div>

        <Button variant="outline" onClick={handleLogout} className="w-full">
          <LogOut className="mr-2 h-4 w-4" />
          Change API Key
        </Button>
      </CardContent>
    </Card>
  );
}
