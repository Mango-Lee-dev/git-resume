import type { Category } from '@/types';

export const CATEGORIES: Category[] = [
  'Feature',
  'Fix',
  'Refactor',
  'Test',
  'Docs',
  'Chore',
];

export const CATEGORY_COLORS: Record<Category, string> = {
  Feature: 'bg-green-100 text-green-800',
  Fix: 'bg-red-100 text-red-800',
  Refactor: 'bg-blue-100 text-blue-800',
  Test: 'bg-purple-100 text-purple-800',
  Docs: 'bg-yellow-100 text-yellow-800',
  Chore: 'bg-gray-100 text-gray-800',
};

export const CATEGORY_ICONS: Record<Category, string> = {
  Feature: '✨',
  Fix: '🐛',
  Refactor: '♻️',
  Test: '🧪',
  Docs: '📝',
  Chore: '🔧',
};
