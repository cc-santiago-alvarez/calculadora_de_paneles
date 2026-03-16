interface IdeamZone {
  department: string;
  capital: string;
  annualAvgGHI: number;
  latitude: number;
  longitude: number;
}

const IDEAM_ZONES: IdeamZone[] = [
  { department: 'La Guajira', capital: 'Riohacha', annualAvgGHI: 6.0, latitude: 11.54, longitude: -72.91 },
  { department: 'Atlántico', capital: 'Barranquilla', annualAvgGHI: 5.5, latitude: 10.96, longitude: -74.78 },
  { department: 'Magdalena', capital: 'Santa Marta', annualAvgGHI: 5.3, latitude: 11.24, longitude: -74.20 },
  { department: 'Cesar', capital: 'Valledupar', annualAvgGHI: 5.2, latitude: 10.47, longitude: -73.25 },
  { department: 'Norte de Santander', capital: 'Cúcuta', annualAvgGHI: 4.8, latitude: 7.89, longitude: -72.50 },
  { department: 'Santander', capital: 'Bucaramanga', annualAvgGHI: 4.6, latitude: 7.12, longitude: -73.12 },
  { department: 'Bolívar', capital: 'Cartagena', annualAvgGHI: 5.0, latitude: 10.39, longitude: -75.51 },
  { department: 'Sucre', capital: 'Sincelejo', annualAvgGHI: 5.1, latitude: 9.30, longitude: -75.39 },
  { department: 'Córdoba', capital: 'Montería', annualAvgGHI: 4.9, latitude: 8.75, longitude: -75.88 },
  { department: 'Antioquia', capital: 'Medellín', annualAvgGHI: 4.5, latitude: 6.25, longitude: -75.56 },
  { department: 'Boyacá', capital: 'Tunja', annualAvgGHI: 4.4, latitude: 5.53, longitude: -73.36 },
  { department: 'Cundinamarca', capital: 'Bogotá', annualAvgGHI: 4.3, latitude: 4.61, longitude: -74.08 },
  { department: 'Tolima', capital: 'Ibagué', annualAvgGHI: 4.7, latitude: 4.44, longitude: -75.24 },
  { department: 'Huila', capital: 'Neiva', annualAvgGHI: 4.8, latitude: 2.93, longitude: -75.28 },
  { department: 'Valle del Cauca', capital: 'Cali', annualAvgGHI: 4.5, latitude: 3.45, longitude: -76.53 },
  { department: 'Cauca', capital: 'Popayán', annualAvgGHI: 4.2, latitude: 2.44, longitude: -76.61 },
  { department: 'Nariño', capital: 'Pasto', annualAvgGHI: 4.0, latitude: 1.21, longitude: -77.28 },
  { department: 'Meta', capital: 'Villavicencio', annualAvgGHI: 4.6, latitude: 4.15, longitude: -73.64 },
  { department: 'Casanare', capital: 'Yopal', annualAvgGHI: 4.8, latitude: 5.34, longitude: -72.39 },
  { department: 'Arauca', capital: 'Arauca', annualAvgGHI: 4.9, latitude: 7.08, longitude: -70.76 },
  { department: 'Risaralda', capital: 'Pereira', annualAvgGHI: 4.3, latitude: 4.81, longitude: -75.69 },
  { department: 'Caldas', capital: 'Manizales', annualAvgGHI: 4.3, latitude: 5.07, longitude: -75.52 },
  { department: 'Quindío', capital: 'Armenia', annualAvgGHI: 4.2, latitude: 4.53, longitude: -75.68 },
  { department: 'Chocó', capital: 'Quibdó', annualAvgGHI: 3.5, latitude: 5.69, longitude: -76.66 },
  { department: 'Amazonas', capital: 'Leticia', annualAvgGHI: 4.0, latitude: -4.22, longitude: -69.94 },
  { department: 'San Andrés', capital: 'San Andrés', annualAvgGHI: 5.6, latitude: 12.58, longitude: -81.70 },
];

export function findNearestZone(lat: number, lng: number): { department: string; capital: string; annualAvgGHI: number } {
  let nearest = IDEAM_ZONES[0];
  let minDist = Infinity;

  for (const zone of IDEAM_ZONES) {
    const dLat = lat - zone.latitude;
    const dLng = lng - zone.longitude;
    const dist = dLat * dLat + dLng * dLng;
    if (dist < minDist) {
      minDist = dist;
      nearest = zone;
    }
  }

  return { department: nearest.department, capital: nearest.capital, annualAvgGHI: nearest.annualAvgGHI };
}
