'use client';

import { CSS_CLASSES } from '../../utils/constants';

interface StatusMessageProps {
  message: string;
  type: 'error' | 'info';
}

/**
 * ステータスメッセージ表示コンポーネント
 */
export default function StatusMessage({ message, type }: StatusMessageProps) {
  const className = type === 'error' ? CSS_CLASSES.ERROR_TEXT : CSS_CLASSES.INFO_TEXT;

  return <p className={className}>{message}</p>;
}