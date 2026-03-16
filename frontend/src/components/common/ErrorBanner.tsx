import { XMarkIcon } from '@heroicons/react/24/outline';

interface ErrorBannerProps {
  message: string;
  onDismiss?: () => void;
}

export default function ErrorBanner({ message, onDismiss }: ErrorBannerProps) {
  if (!message) return null;

  return (
    <div className="bg-danger-soft border border-[var(--color-danger)] border-opacity-20 text-danger px-4 py-3 rounded-lg flex items-center justify-between mb-4">
      <p className="text-sm">{message}</p>
      {onDismiss && (
        <button
          onClick={onDismiss}
          className="text-danger hover:text-[hsl(0,72%,40%)] ml-4 transition-colors duration-fast"
        >
          <XMarkIcon className="w-4 h-4" />
        </button>
      )}
    </div>
  );
}
