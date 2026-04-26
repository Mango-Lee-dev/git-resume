'use client';

import { Badge } from '@/components/ui/badge';
import { CATEGORY_COLORS, CATEGORY_ICONS } from '@/lib/constants';
import type { Category } from '@/types';

interface CategoryBadgeProps {
  category: Category;
}

export function CategoryBadge({ category }: CategoryBadgeProps) {
  const icon = CATEGORY_ICONS[category] || '';
  const colorClass = CATEGORY_COLORS[category] || 'bg-gray-100 text-gray-800';

  return (
    <Badge variant="secondary" className={colorClass}>
      {icon} {category}
    </Badge>
  );
}
