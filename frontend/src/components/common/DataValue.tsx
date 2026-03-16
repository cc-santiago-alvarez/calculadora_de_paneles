interface DataValueProps {
  value: string | number;
  unit?: string;
  size?: 'sm' | 'md' | 'lg';
}

const sizeMap = {
  sm: 'text-sm',
  md: 'text-lg',
  lg: 'text-2xl',
};

export default function DataValue({ value, unit, size = 'md' }: DataValueProps) {
  return (
    <span className={`font-mono tabular-nums text-fg-primary font-semibold ${sizeMap[size]}`}>
      {value}
      {unit && <span className="text-fg-tertiary font-normal text-[0.75em] ml-1">{unit}</span>}
    </span>
  );
}
