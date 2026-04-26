'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { X, FolderGit2 } from 'lucide-react';
import { useAnalyzeStore } from '@/stores';

export function RepoStep() {
  const { repos, addRepo, removeRepo } = useAnalyzeStore();
  const [inputValue, setInputValue] = useState('');

  const handleAdd = () => {
    if (inputValue.trim()) {
      addRepo(inputValue.trim());
      setInputValue('');
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAdd();
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <FolderGit2 className="h-5 w-5" />
          Select Repositories
        </CardTitle>
        <CardDescription>
          Enter the paths to the Git repositories you want to analyze
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex gap-2">
          <Input
            placeholder="/path/to/repository"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={handleKeyDown}
          />
          <Button onClick={handleAdd} disabled={!inputValue.trim()}>
            Add
          </Button>
        </div>

        {repos.length > 0 && (
          <div className="space-y-2">
            {repos.map((repo, index) => (
              <div
                key={index}
                className="flex items-center justify-between p-3 bg-muted rounded-md"
              >
                <span className="text-sm font-mono truncate">{repo}</span>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => removeRepo(index)}
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>
            ))}
          </div>
        )}

        {repos.length === 0 && (
          <p className="text-sm text-muted-foreground text-center py-4">
            No repositories added yet
          </p>
        )}
      </CardContent>
    </Card>
  );
}
