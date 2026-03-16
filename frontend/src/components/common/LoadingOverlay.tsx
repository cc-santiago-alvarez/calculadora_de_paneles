export default function LoadingOverlay({ message = 'Cargando...' }: { message?: string }) {
  return (
    <div className="fixed inset-0 bg-black/20 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-surface-raised rounded-lg p-8 border border-[var(--color-border-default)] flex flex-col items-center gap-4">
        <svg className="w-10 h-10 animate-spin" viewBox="0 0 40 40" fill="none">
          <circle cx="20" cy="20" r="17" stroke="var(--color-border-default)" strokeWidth="3" />
          <path
            d="M37 20a17 17 0 0 0-17-17"
            stroke="var(--color-brand)"
            strokeWidth="3"
            strokeLinecap="round"
          />
        </svg>
        <p className="text-fg-secondary font-medium text-sm">{message}</p>
      </div>
    </div>
  );
}
