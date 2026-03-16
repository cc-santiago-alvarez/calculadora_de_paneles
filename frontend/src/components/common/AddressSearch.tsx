import { useState, useRef, useEffect } from 'react';
import { MagnifyingGlassIcon } from '@heroicons/react/24/outline';
import { geocodeApi, type GeoSearchResult } from '../../api/geocode';

interface AddressSearchProps {
  onSelect: (result: GeoSearchResult) => void;
}

export default function AddressSearch({ onSelect }: AddressSearchProps) {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<GeoSearchResult[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleChange = (value: string) => {
    setQuery(value);

    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    if (value.length < 3) {
      setResults([]);
      setIsOpen(false);
      return;
    }

    debounceRef.current = setTimeout(async () => {
      setIsLoading(true);
      try {
        const data = await geocodeApi.searchAddress(value);
        setResults(data);
        setIsOpen(data.length > 0);
      } catch {
        setResults([]);
        setIsOpen(false);
      } finally {
        setIsLoading(false);
      }
    }, 400);
  };

  const handleSelect = (result: GeoSearchResult) => {
    setQuery(result.displayName);
    setIsOpen(false);
    setResults([]);
    onSelect(result);
  };

  return (
    <div ref={containerRef} className="relative">
      <div className="relative">
        <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted pointer-events-none" />
        <input
          type="text"
          value={query}
          onChange={(e) => handleChange(e.target.value)}
          onFocus={() => results.length > 0 && setIsOpen(true)}
          placeholder="Buscar direccion o ciudad..."
          className="w-full bg-inset border border-[var(--color-border-default)] rounded-md py-2 pl-9 pr-9 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2"
        />
        {isLoading && (
          <div className="absolute right-3 top-1/2 -translate-y-1/2">
            <div className="w-4 h-4 border-2 border-brand border-t-transparent rounded-full animate-spin" />
          </div>
        )}
      </div>

      {isOpen && results.length > 0 && (
        <ul className="absolute z-20 w-full mt-1 bg-surface-raised border border-[var(--color-border-default)] rounded-lg shadow-lg overflow-hidden animate-fade-in-up">
          {results.map((result, i) => (
            <li key={i}>
              <button
                type="button"
                onClick={() => handleSelect(result)}
                className="w-full text-left px-4 py-2.5 text-sm text-fg-primary hover:bg-inset transition-colors duration-fast border-b border-[var(--color-border-subtle)] last:border-0"
              >
                <span className="line-clamp-2">{result.displayName}</span>
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
