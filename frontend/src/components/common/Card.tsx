import { type ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  padding?: 'none' | 'sm' | 'md' | 'lg';
  hoverable?: boolean;
  className?: string;
  onClick?: () => void;
}

const paddingMap = {
  none: '',
  sm: 'p-4',
  md: 'p-5',
  lg: 'p-6',
};

export default function Card({ children, padding = 'md', hoverable = false, className = '', onClick }: CardProps) {
  return (
    <div
      onClick={onClick}
      className={`bg-surface rounded-lg border border-[var(--color-border-default)] animate-fade-in-up ${paddingMap[padding]} ${
        hoverable
          ? 'transition-all duration-fast ease-decel hover:border-[var(--color-border-strong)] hover:shadow-md hover:-translate-y-0.5 cursor-pointer'
          : ''
      } ${className}`}
    >
      {children}
    </div>
  );
}
