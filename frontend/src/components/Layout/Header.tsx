import { useProjectStore } from '../../store/useProjectStore';
import ThemeToggle from '../common/ThemeToggle';

export default function Header() {
  const { currentProject, isLoading } = useProjectStore();

  return (
    <header className="bg-surface border-b border-[var(--color-border-subtle)] px-6 py-3 flex items-center justify-between">
      <div>
        {currentProject ? (
          <div>
            <h2 className="text-base font-semibold text-fg-primary">{currentProject.name}</h2>
            <p className="text-xs text-fg-tertiary">
              {currentProject.location.city}, {currentProject.location.department}
            </p>
          </div>
        ) : (
          <h2 className="text-base font-semibold text-fg-primary">Calculadora de Paneles Solares</h2>
        )}
      </div>
      <div className="flex items-center gap-4">
        {isLoading && (
          <div className="flex items-center gap-2 text-sm text-fg-tertiary">
            <div className="w-4 h-4 border-2 border-brand border-t-transparent rounded-full animate-spin" />
            Calculando...
          </div>
        )}
        <ThemeToggle />
      </div>
    </header>
  );
}
