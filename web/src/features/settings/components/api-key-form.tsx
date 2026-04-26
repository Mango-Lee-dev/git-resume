'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Key, Eye, EyeOff, CheckCircle2, AlertCircle } from 'lucide-react';
import { useSessionStore } from '@/stores';
import { useCreateSession } from '../hooks/use-create-session';

export function ApiKeyForm() {
  const router = useRouter();
  const [apiKey, setApiKey] = useState('');
  const [showKey, setShowKey] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const { setApiKey: storeApiKey, setSession } = useSessionStore();
  const { mutate: createSession, isPending } = useCreateSession();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!apiKey.trim()) {
      setError('Please enter your API key');
      return;
    }

    if (!apiKey.startsWith('sk-ant-')) {
      setError('Invalid API key format. Should start with sk-ant-');
      return;
    }

    createSession(
      { api_key: apiKey },
      {
        onSuccess: (data) => {
          storeApiKey(apiKey);
          setSession({
            id: data.session_id,
            createdAt: data.created_at,
            expiresAt: data.expires_at,
          });
          router.push('/');
        },
        onError: (err) => {
          setError(err instanceof Error ? err.message : 'Failed to validate API key');
        },
      }
    );
  };

  return (
    <Card className="max-w-md mx-auto">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Key className="h-5 w-5" />
          Claude API Key
        </CardTitle>
        <CardDescription>
          Enter your Anthropic API key to start analyzing commits.
          Your key is stored locally and sent securely with each request.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">API Key</label>
            <div className="relative">
              <Input
                type={showKey ? 'text' : 'password'}
                placeholder="sk-ant-api03-..."
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
                className="pr-10"
              />
              <button
                type="button"
                onClick={() => setShowKey(!showKey)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              >
                {showKey ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </button>
            </div>
          </div>

          {error && (
            <div className="flex items-center gap-2 text-sm text-destructive">
              <AlertCircle className="h-4 w-4" />
              {error}
            </div>
          )}

          <Button type="submit" className="w-full" disabled={isPending}>
            {isPending ? (
              'Validating...'
            ) : (
              <>
                <CheckCircle2 className="mr-2 h-4 w-4" />
                Validate & Continue
              </>
            )}
          </Button>

          <p className="text-xs text-muted-foreground text-center">
            Get your API key from{' '}
            <a
              href="https://console.anthropic.com/settings/keys"
              target="_blank"
              rel="noopener noreferrer"
              className="underline hover:text-foreground"
            >
              console.anthropic.com
            </a>
          </p>
        </form>
      </CardContent>
    </Card>
  );
}
