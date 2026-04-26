'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useSessionStore } from '@/stores';

const PUBLIC_PATHS = ['/settings'];

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const { session, isAuthenticated } = useSessionStore();
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    // Skip auth check for public paths
    if (PUBLIC_PATHS.includes(pathname)) {
      setIsChecking(false);
      return;
    }

    // Check if authenticated
    if (!isAuthenticated()) {
      router.replace('/settings');
    } else {
      setIsChecking(false);
    }
  }, [pathname, session, isAuthenticated, router]);

  // Show nothing while checking auth on protected routes
  if (isChecking && !PUBLIC_PATHS.includes(pathname)) {
    return (
      <div className="flex items-center justify-center min-h-[50vh]">
        <div className="text-muted-foreground">Loading...</div>
      </div>
    );
  }

  return <>{children}</>;
}
