import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/react';
import { CheckIcon, ChevronUpDownIcon } from '@heroicons/react/24/outline';

export interface SelectOption {
  value: string;
  label: string;
}

interface SelectProps {
  value: string;
  onChange: (value: string) => void;
  options: SelectOption[];
  placeholder?: string;
  disabled?: boolean;
  className?: string;
}

export default function Select({ value, onChange, options, placeholder = 'Seleccionar...', disabled = false, className = '' }: SelectProps) {
  const selected = options.find((o) => o.value === value);

  return (
    <Listbox value={value} onChange={onChange} disabled={disabled}>
      <div className={`relative ${className}`}>
        <ListboxButton
          className={`relative w-full rounded-md bg-inset border border-[var(--color-border-default)] py-2 pl-3 pr-10 text-left text-sm transition-colors duration-fast ease-decel focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2 hover:border-[var(--color-border-strong)] disabled:opacity-50 ${
            selected ? 'text-fg-primary' : 'text-fg-muted'
          }`}
        >
          <span className="block truncate">{selected ? selected.label : placeholder}</span>
          <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
            <ChevronUpDownIcon className="h-4 w-4 text-fg-muted" />
          </span>
        </ListboxButton>

        <ListboxOptions
          className="absolute z-20 mt-1 max-h-60 w-full overflow-auto rounded-md bg-surface-raised border border-[var(--color-border-default)] py-1 text-sm shadow-lg focus:outline-none"
        >
          {options.map((option) => (
            <ListboxOption
              key={option.value}
              value={option.value}
              className={({ active, selected: sel }) =>
                `relative cursor-pointer select-none py-2 pl-3 pr-9 transition-colors duration-fast ${
                  active ? 'bg-brand-soft text-fg-primary' : 'text-fg-secondary'
                } ${sel ? 'font-medium text-fg-primary' : ''}`
              }
            >
              {({ selected: sel }) => (
                <>
                  <span className="block truncate">{option.label}</span>
                  {sel && (
                    <span className="absolute inset-y-0 right-0 flex items-center pr-3 text-brand">
                      <CheckIcon className="h-4 w-4" />
                    </span>
                  )}
                </>
              )}
            </ListboxOption>
          ))}
        </ListboxOptions>
      </div>
    </Listbox>
  );
}
