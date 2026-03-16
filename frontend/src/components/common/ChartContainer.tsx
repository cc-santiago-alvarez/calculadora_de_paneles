import { type ReactNode } from 'react';
import Card from './Card';

interface ChartContainerProps {
  title: string;
  children: ReactNode;
  className?: string;
}

export default function ChartContainer({ title, children, className = '' }: ChartContainerProps) {
  return (
    <Card padding="none" className={className}>
      <div className="px-5 pt-5 pb-2">
        <h3 className="text-sm font-medium text-fg-secondary">{title}</h3>
      </div>
      <div className="px-3 pb-4">{children}</div>
    </Card>
  );
}
