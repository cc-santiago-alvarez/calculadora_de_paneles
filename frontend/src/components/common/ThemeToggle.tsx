import { useState } from 'react';
import { SunIcon, MoonIcon } from '@heroicons/react/24/outline';
import { useTheme } from '../../hooks/useTheme';

export default function ThemeToggle() {
  const { theme, toggleTheme } = useTheme();
  const [animating, setAnimating] = useState(false);

  const handleClick = () => {
    setAnimating(true);
    toggleTheme();
    setTimeout(() => setAnimating(false), 300);
  };

  return (
    <button
      onClick={handleClick}
      className="relative p-2 rounded-lg text-fg-secondary hover:text-fg-primary hover:bg-inset transition-colors duration-fast focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2"
      aria-label={theme === 'light' ? 'Activar modo oscuro' : 'Activar modo claro'}
    >
      {theme === 'light' ? (
        <MoonIcon className={`w-5 h-5 ${animating ? 'animate-theme-icon' : ''}`} />
      ) : (
        <SunIcon className={`w-5 h-5 ${animating ? 'animate-theme-icon' : ''}`} />
      )}
    </button>
  );
}
