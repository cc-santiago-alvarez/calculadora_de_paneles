export const PANEL_FORMATS = [
  {
    key: 'standard' as const,
    label: 'Estandar',
    watts: 400,
    area: 1.95,
    minPower: 350,
    maxPower: 429,
    description: 'Paneles residenciales compactos, ideales para techos pequenos',
  },
  {
    key: 'high-efficiency' as const,
    label: 'Alta eficiencia',
    watts: 450,
    area: 2.05,
    minPower: 430,
    maxPower: 529,
    description: 'Mayor potencia por panel, buen equilibrio entre tamano y rendimiento',
  },
  {
    key: 'large-format' as const,
    label: 'Gran formato',
    watts: 550,
    area: 2.58,
    minPower: 530,
    maxPower: 999,
    description: 'Paneles de mayor potencia, requieren mas espacio pero menos unidades',
  },
] as const;
