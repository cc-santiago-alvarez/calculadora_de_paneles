import api from './client';

export interface GeoSearchResult {
  displayName: string;
  lat: number;
  lon: number;
}

export const geocodeApi = {
  reverseGeocode: (lat: number, lon: number) =>
    api.get<{ city: string; department: string }>(`/geocode/reverse?lat=${lat}&lon=${lon}`).then((r) => r.data),

  searchAddress: (query: string) =>
    api.get<GeoSearchResult[]>('/geocode/search', { params: { q: query } }).then((r) => r.data),
};
