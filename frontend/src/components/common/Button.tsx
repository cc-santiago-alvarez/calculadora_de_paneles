import { type ButtonHTMLAttributes, type ReactNode } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  children: ReactNode;
}

const base =
  'inline-flex items-center justify-center font-medium transition-all duration-fast ease-decel active:scale-[0.98] focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2 disabled:opacity-50 disabled:pointer-events-none';

const variants = {
  primary:
    'bg-brand text-white hover:bg-brand-hover shadow-sm hover:shadow-md rounded-md',
  secondary:
    'bg-transparent border border-[var(--color-border-default)] text-fg-primary hover:border-[var(--color-border-strong)] hover:bg-inset rounded-md',
  ghost:
    'bg-transparent text-fg-secondary hover:bg-inset hover:text-fg-primary rounded-md',
  danger:
    'bg-danger text-white hover:brightness-90 rounded-md',
};

const sizes = {
  sm: 'text-xs px-3 py-1.5 gap-1.5',
  md: 'text-sm px-4 py-2 gap-2',
  lg: 'text-sm px-5 py-2.5 gap-2',
};

export default function Button({ variant = 'primary', size = 'md', className = '', children, ...props }: ButtonProps) {
  return (
    <button className={`${base} ${variants[variant]} ${sizes[size]} ${className}`} {...props}>
      {children}
    </button>
  );
}
