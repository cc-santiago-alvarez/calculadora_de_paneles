import { useProjectStore } from '../../store/useProjectStore';
import { MapContainer, TileLayer, Marker, useMapEvents, useMap } from 'react-leaflet';
import { useRef, useCallback, useState, useEffect } from 'react';
import { findNearestZone } from '../../utils/ideamZones';
import { calculationApi } from '../../api/calculation';
import { geocodeApi, type GeoSearchResult } from '../../api/geocode';
import Select from '../../components/common/Select';
import AddressSearch from '../../components/common/AddressSearch';

function useLocationUpdater() {
  const { formData, setFormData, setIrradiationPreview } = useProjectStore();
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const latestRef = useRef<{ lat: number; lng: number } | null>(null);

  const updateLocation = useCallback(
    (lat: number, lng: number) => {
      const nearest = findNearestZone(lat, lng);
      setFormData({
        location: {
          ...formData.location,
          latitude: lat,
          longitude: lng,
          department: nearest.department,
          city: nearest.capital,
        },
      });

      latestRef.current = { lat, lng };

      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }

      debounceRef.current = setTimeout(async () => {
        const point = { lat, lng };

        try {
          const irradiation = await calculationApi.fetchIrradiation({
            latitude: lat,
            longitude: lng,
            tilt: formData.roof.slopes.length > 0 ? formData.roof.slopes[0].tilt : formData.roof.tilt,
            azimuth: formData.roof.slopes.length > 0 ? formData.roof.slopes[0].azimuth : formData.roof.azimuth,
          });

          if (latestRef.current?.lat === point.lat && latestRef.current?.lng === point.lng) {
            setIrradiationPreview(irradiation);
            setFormData({
              location: {
                ...useProjectStore.getState().formData.location,
                altitude: Math.round(irradiation.elevation),
              },
            });
          }
        } catch {
          // Silently fail
        }

        try {
          const geo = await geocodeApi.reverseGeocode(lat, lng);
          if (latestRef.current?.lat === point.lat && latestRef.current?.lng === point.lng) {
            const currentLocation = useProjectStore.getState().formData.location;
            const updates: Partial<typeof currentLocation> = {};
            if (geo.city) updates.city = geo.city;
            if (geo.department) updates.department = geo.department;
            if (Object.keys(updates).length > 0) {
              setFormData({
                location: { ...useProjectStore.getState().formData.location, ...updates },
              });
            }
          }
        } catch {
          // Silently fail
        }
      }, 300);
    },
    [formData.location, formData.roof.tilt, formData.roof.azimuth, formData.roof.slopes, setFormData, setIrradiationPreview],
  );

  return updateLocation;
}

function MapFlyTo({ target }: { target: { lat: number; lng: number } | null }) {
  const map = useMap();
  useEffect(() => {
    if (target) {
      map.flyTo([target.lat, target.lng], 13, { duration: 1 });
    }
  }, [target, map]);
  return null;
}

function LocationPicker({ onLocationChange }: { onLocationChange: (lat: number, lng: number) => void }) {
  const { formData } = useProjectStore();

  useMapEvents({
    click(e) {
      onLocationChange(
        parseFloat(e.latlng.lat.toFixed(4)),
        parseFloat(e.latlng.lng.toFixed(4)),
      );
    },
  });

  return (
    <Marker position={[formData.location.latitude, formData.location.longitude]} />
  );
}

const DEPARTMENTS = [
  'Amazonas', 'Antioquia', 'Arauca', 'Atlántico', 'Bolívar', 'Boyacá',
  'Caldas', 'Casanare', 'Cauca', 'Cesar', 'Chocó', 'Córdoba',
  'Cundinamarca', 'Huila', 'La Guajira', 'Magdalena', 'Meta',
  'Nariño', 'Norte de Santander', 'Quindío', 'Risaralda',
  'San Andrés', 'Santander', 'Sucre', 'Tolima', 'Valle del Cauca',
];

const DEPARTMENT_OPTIONS = DEPARTMENTS.map((d) => ({
  value: d,
  label: d,
}));

export default function LocationStep() {
  const { formData, setFormData } = useProjectStore();
  const updateLocation = useLocationUpdater();
  const [flyTarget, setFlyTarget] = useState<{ lat: number; lng: number } | null>(null);

  const handleSearchSelect = (result: GeoSearchResult) => {
    const lat = parseFloat(result.lat.toFixed(4));
    const lng = parseFloat(result.lon.toFixed(4));
    setFlyTarget({ lat, lng });
    updateLocation(lat, lng);
  };

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-fg-primary mb-2">Ubicacion del Proyecto</h3>
        <p className="text-sm text-fg-muted mb-4">
          Busca una direccion, haz clic en el mapa o ingresa las coordenadas manualmente.
        </p>
      </div>

      {/* Address search */}
      <AddressSearch onSelect={handleSearchSelect} />

      {/* Map */}
      <div className="h-[350px] rounded-lg overflow-hidden border border-[var(--color-border-default)]">
        <MapContainer
          center={[formData.location.latitude, formData.location.longitude]}
          zoom={6}
          className="h-full w-full"
        >
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <LocationPicker onLocationChange={updateLocation} />
          <MapFlyTo target={flyTarget} />
        </MapContainer>
      </div>

      {/* Coordinate inputs */}
      <div className="grid grid-cols-3 gap-4">
        <div>
          <label className="block text-sm font-medium text-fg-secondary mb-1">Latitud</label>
          <input
            type="number"
            step="0.0001"
            value={formData.location.latitude}
            onChange={(e) =>
              setFormData({ location: { ...formData.location, latitude: parseFloat(e.target.value) || 0 } })
            }
            className="w-full bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2 font-mono tabular-nums"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-fg-secondary mb-1">Longitud</label>
          <input
            type="number"
            step="0.0001"
            value={formData.location.longitude}
            onChange={(e) =>
              setFormData({ location: { ...formData.location, longitude: parseFloat(e.target.value) || 0 } })
            }
            className="w-full bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2 font-mono tabular-nums"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-fg-secondary mb-1">Altitud (m)</label>
          <input
            type="number"
            value={formData.location.altitude}
            onChange={(e) =>
              setFormData({ location: { ...formData.location, altitude: parseFloat(e.target.value) || 0 } })
            }
            className="w-full bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2 font-mono tabular-nums"
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-fg-secondary mb-1">Departamento</label>
          <Select
            value={formData.location.department}
            onChange={(val) =>
              setFormData({ location: { ...formData.location, department: val } })
            }
            options={DEPARTMENT_OPTIONS}
            placeholder="Seleccionar..."
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-fg-secondary mb-1">Ciudad</label>
          <input
            type="text"
            value={formData.location.city}
            onChange={(e) =>
              setFormData({ location: { ...formData.location, city: e.target.value } })
            }
            className="w-full bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2"
          />
        </div>
      </div>

    </div>
  );
}
