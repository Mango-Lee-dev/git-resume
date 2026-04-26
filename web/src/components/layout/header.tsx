'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { GitBranch, LayoutDashboard, PlayCircle, FileText, Download, Settings, Key } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useSessionStore } from '@/stores';

const navItems = [
  { href: '/', label: 'Dashboard', icon: LayoutDashboard },
  { href: '/analyze', label: 'Analyze', icon: PlayCircle },
  { href: '/results', label: 'Results', icon: FileText },
  { href: '/export', label: 'Export', icon: Download },
];

export function Header() {
  const pathname = usePathname();
  const { session, apiKey } = useSessionStore();
  const isAuthenticated = !!session && !!apiKey;

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 items-center">
        <Link href="/" className="flex items-center gap-2 font-semibold mr-8">
          <GitBranch className="h-5 w-5" />
          <span>Git Resume</span>
        </Link>

        <nav className="flex items-center gap-1 flex-1">
          {navItems.map((item) => {
            const isActive =
              pathname === item.href ||
              (item.href !== '/' && pathname.startsWith(item.href));

            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  'flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-md transition-colors',
                  isActive
                    ? 'bg-accent text-accent-foreground'
                    : 'text-muted-foreground hover:text-foreground hover:bg-accent/50'
                )}
              >
                <item.icon className="h-4 w-4" />
                {item.label}
              </Link>
            );
          })}
        </nav>

        <Link
          href="/settings"
          className={cn(
            'flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-md transition-colors',
            pathname === '/settings'
              ? 'bg-accent text-accent-foreground'
              : 'text-muted-foreground hover:text-foreground hover:bg-accent/50'
          )}
        >
          {isAuthenticated ? (
            <Settings className="h-4 w-4" />
          ) : (
            <Key className="h-4 w-4" />
          )}
          <span className="hidden sm:inline">
            {isAuthenticated ? 'Settings' : 'API Key'}
          </span>
        </Link>
      </div>
    </header>
  );
}
