import { NavLink } from 'react-router-dom';
import logoLight from '../../assets/images/logo-codecraft.png';
import logoDark from '../../assets/images/logo-codecraft_blanco.png';
import {
  FolderIcon,
  PlusCircleIcon,
  MapIcon,
  ChartBarIcon,
  ScaleIcon,
  DocumentTextIcon,
} from '@heroicons/react/24/outline';

const navItems = [
  { path: '/', label: 'Proyectos', icon: FolderIcon },
  { path: '/new', label: 'Nuevo Proyecto', icon: PlusCircleIcon },
  { path: '/map', label: 'Mapa', icon: MapIcon },
  { path: '/results', label: 'Resultados', icon: ChartBarIcon },
  { path: '/compare', label: 'Comparador', icon: ScaleIcon },
  { path: '/report', label: 'Reportes', icon: DocumentTextIcon },
];

export default function Sidebar() {
  return (
    <aside className="w-64 bg-canvas flex flex-col min-h-screen border-r border-[var(--color-border-default)]">
      <div className="p-5 border-b border-[var(--color-border-default)]">
        <h1 className="text-base font-semibold bg-gradient-to-r from-brand to-solar-600 bg-clip-text text-transparent">
          Solar Panel Calculator
        </h1>
        <p className="text-xs text-fg-tertiary mt-1 tracking-wide uppercase">
          Dimensionamiento FV Colombia
        </p>
      </div>
      <nav className="flex-1 py-3">
        {navItems.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            className={({ isActive }) =>
              `group flex items-center gap-3 px-4 py-2.5 text-sm transition-all duration-fast ease-decel ${
                isActive
                  ? 'bg-brand-soft text-brand border-l-2 border-brand font-medium'
                  : 'text-fg-secondary hover:bg-inset hover:text-fg-primary'
              }`
            }
          >
            <item.icon className="w-5 h-5 flex-shrink-0 transition-transform duration-fast group-hover:scale-110" />
            <span>{item.label}</span>
          </NavLink>
        ))}
      </nav>
      <div className="p-4 border-t border-[var(--color-border-default)] flex justify-center">
        <img src={logoLight} alt="Code Craft" className="h-12 w-auto opacity-70 dark:hidden" />
        <img src={logoDark} alt="Code Craft" className="h-12 w-auto opacity-70 hidden dark:block" />
      </div>
    </aside>
  );
}
