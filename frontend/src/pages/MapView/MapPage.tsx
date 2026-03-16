import { useState } from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMapEvents } from 'react-leaflet';
import { useProjectStore } from '../../store/useProjectStore';
import { calculationApi } from '../../api/calculation';
import { IrradiationResult } from '../../types';

function MapClickHandler({ onLocationSelect }: { onLocationSelect: (lat: number, lon: number) => void }) {
  useMapEvents({
    click(e) {
      onLocationSelect(parseFloat(e.latlng.lat.toFixed(4)), parseFloat(e.latlng.lng.toFixed(4)));
    },
  });
  return null;
}

export default function MapPage() {
  const { setIrradiationPreview, irradiationPreview } = useProjectStore();
  const [selectedPos, setSelectedPos] = useState<{ lat: number; lon: number } | null>(null);
  const [loading, setLoading] = useState(false);

  const handleLocationSelect = async (lat: number, lon: number) => {
    setSelectedPos({ lat, lon });
    setLoading(true);
    try {
      const result = await calculationApi.fetchIrradiation({
        latitude: lat,
        longitude: lon,
        tilt: 10,
        azimuth: 0,
      });
      setIrradiationPreview(result);
    } catch (err) {
      console.error('Failed to fetch irradiation:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="h-full flex gap-4">
      <div className="flex-1 rounded-lg overflow-hidden border border-[var(--color-border-default)]">
        <MapContainer center={[4.61, -74.08]} zoom={6} className="h-full w-full">
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <MapClickHandler onLocationSelect={handleLocationSelect} />
          {selectedPos && (
            <Marker position={[selectedPos.lat, selectedPos.lon]}>
              <Popup>
                <div className="text-sm">
                  <p className="font-semibold">
                    {selectedPos.lat}, {selectedPos.lon}
                  </p>
                  {irradiationPreview && (
                    <p className="mt-1">
                      HSP: {irradiationPreview.annualAvgHSP.toFixed(2)} h/dia
                    </p>
                  )}
                </div>
              </Popup>
            </Marker>
          )}
        </MapContainer>
      </div>

      {/* Side panel */}
      <div className="w-80 bg-surface rounded-lg border border-[var(--color-border-default)] p-5 overflow-auto">
        <h3 className="font-semibold text-fg-primary mb-4">Consulta de Irradiacion</h3>
        <p className="text-sm text-fg-secondary mb-4">
          Haz clic en el mapa para consultar la irradiacion solar en cualquier punto de Colombia.
        </p>

        {loading && (
          <div className="flex items-center gap-2 text-sm text-fg-secondary mb-4">
            <div className="w-4 h-4 border-2 border-brand border-t-transparent rounded-full animate-spin" />
            Consultando datos...
          </div>
        )}

        {selectedPos && (
          <div className="mb-4 bg-inset rounded-md p-3">
            <p className="text-xs text-fg-tertiary">Coordenadas seleccionadas</p>
            <p className="font-mono tabular-nums text-sm text-fg-primary">
              {selectedPos.lat}, {selectedPos.lon}
            </p>
          </div>
        )}

        {irradiationPreview && (
          <div className="space-y-4">
            <div className="bg-brand-soft rounded-lg p-4">
              <p className="text-xs text-fg-tertiary uppercase">HSP Promedio Anual</p>
              <p className="text-2xl font-bold font-mono tabular-nums text-brand">
                {irradiationPreview.annualAvgHSP.toFixed(2)}
              </p>
              <p className="text-xs text-fg-tertiary">kWh/m2/dia</p>
            </div>

            <div className="bg-inset rounded-md p-3">
              <p className="text-xs text-fg-tertiary mb-1">Fuente: {irradiationPreview.source}</p>
            </div>

            <div>
              <p className="text-sm font-medium text-fg-primary mb-2">GHI Mensual (kWh/m2/dia)</p>
              <div className="space-y-1">
                {['Ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul', 'Ago', 'Sep', 'Oct', 'Nov', 'Dic'].map(
                  (month, i) => (
                    <div key={month} className="flex items-center gap-2 text-sm">
                      <span className="w-10 text-fg-tertiary">{month}</span>
                      <div className="flex-1 bg-[var(--color-border-default)] rounded-full h-2">
                        <div
                          className="bg-brand rounded-full h-2"
                          style={{
                            width: `${(irradiationPreview.monthlyGHI[i] / 7) * 100}%`,
                          }}
                        />
                      </div>
                      <span className="w-10 text-right font-mono tabular-nums text-xs text-fg-primary">
                        {irradiationPreview.monthlyGHI[i]?.toFixed(1)}
                      </span>
                    </div>
                  )
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
