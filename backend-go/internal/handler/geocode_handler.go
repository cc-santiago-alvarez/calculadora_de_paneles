package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dev13/calculadora-paneles-backend/internal/middleware"
)

// cleanColombianAddress normalizes a Colombian address for geocoding.
// Extracts street type + number (e.g., "Calle 53A") and discards house numbering (# 48-14).
var reColAddr = regexp.MustCompile(`(?i)^(calle|carrera|cra|cll|cr|cl|transversal|tv|diagonal|dg|avenida|av)\s*\.?\s*(\d+\w*)\b`)
var reHashNum = regexp.MustCompile(`#.*$`)

func cleanColombianAddress(query string) string {
	// Remove # and everything after it (house numbering)
	cleaned := reHashNum.ReplaceAllString(query, "")
	cleaned = strings.TrimSpace(cleaned)
	return cleaned
}

func extractStreetName(query string) string {
	m := reColAddr.FindStringSubmatch(query)
	if m == nil {
		return ""
	}
	return m[1] + " " + m[2]
}

type GeocodingHandler struct {
	mu       sync.Mutex
	lastCall time.Time
	client   *http.Client
}

func NewGeocodingHandler() *GeocodingHandler {
	return &GeocodingHandler{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type photonResponse struct {
	Features []photonFeature `json:"features"`
}

type photonFeature struct {
	Properties photonProperties `json:"properties"`
	Geometry   photonGeometry   `json:"geometry"`
}

type photonProperties struct {
	Name        string `json:"name"`
	Street      string `json:"street"`
	HouseNumber string `json:"housenumber"`
	City        string `json:"city"`
	County      string `json:"county"`
	State       string `json:"state"`
	Country     string `json:"country"`
}

type photonGeometry struct {
	Coordinates [2]float64 `json:"coordinates"` // [lon, lat]
}

type searchResult struct {
	DisplayName string  `json:"displayName"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

type nominatimAddress struct {
	State   string `json:"state"`
	City    string `json:"city"`
	Town    string `json:"town"`
	Village string `json:"village"`
	County  string `json:"county"`
}

type nominatimResponse struct {
	Address nominatimAddress `json:"address"`
}

func (h *GeocodingHandler) ReverseGeocode(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	if latStr == "" || lonStr == "" {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Parámetros lat y lon son requeridos",
		})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil || lat < -90 || lat > 90 {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Latitud inválida",
		})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil || lon < -180 || lon > 180 {
		middleware.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Longitud inválida",
		})
		return
	}

	// Rate limiting: wait if less than 1 second since last call
	h.mu.Lock()
	elapsed := time.Since(h.lastCall)
	if elapsed < time.Second {
		time.Sleep(time.Second - elapsed)
	}
	h.lastCall = time.Now()
	h.mu.Unlock()

	url := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json&accept-language=es",
		lat, lon,
	)

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusInternalServerError, "Error creando request"))
		return
	}
	req.Header.Set("User-Agent", "CalculadoraPanelesSolares/1.0")

	resp, err := h.client.Do(req)
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadGateway, "Error consultando servicio de geocodificación"))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadGateway, "Error leyendo respuesta de geocodificación"))
		return
	}

	var nomResp nominatimResponse
	if err := json.Unmarshal(body, &nomResp); err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadGateway, "Error procesando respuesta de geocodificación"))
		return
	}

	city := nomResp.Address.City
	if city == "" {
		city = nomResp.Address.Town
	}
	if city == "" {
		city = nomResp.Address.Village
	}
	if city == "" {
		city = nomResp.Address.County
	}

	department := nomResp.Address.State

	middleware.WriteJSON(w, http.StatusOK, map[string]string{
		"city":       city,
		"department": department,
	})
}

// nominatimSearchResult for Nominatim forward geocoding
type nominatimSearchResult struct {
	DisplayName string `json:"display_name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
}

func (h *GeocodingHandler) SearchAddress(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" || len(query) < 3 {
		middleware.WriteJSON(w, http.StatusOK, []searchResult{})
		return
	}

	// Clean Colombian address: remove # house numbering
	cleaned := cleanColombianAddress(query)

	// Run Photon and Nominatim in parallel
	type apiResult struct {
		results []searchResult
	}

	photonCh := make(chan apiResult, 1)
	nominatimCh := make(chan apiResult, 1)

	// Photon search (no rate limiting needed)
	go func() {
		res := h.searchPhoton(r, cleaned+" Colombia")
		photonCh <- apiResult{results: res}
	}()

	// Nominatim search (with rate limiting)
	go func() {
		res := h.searchNominatim(r, cleaned)
		nominatimCh <- apiResult{results: res}
	}()

	photonRes := <-photonCh
	nominatimRes := <-nominatimCh

	// Merge results: Photon first, then Nominatim (deduplicated)
	seen := map[string]bool{}
	results := make([]searchResult, 0, 8)
	for _, sr := range photonRes.results {
		key := fmt.Sprintf("%.3f,%.3f", sr.Lat, sr.Lon)
		if !seen[key] {
			seen[key] = true
			results = append(results, sr)
		}
	}
	for _, sr := range nominatimRes.results {
		key := fmt.Sprintf("%.3f,%.3f", sr.Lat, sr.Lon)
		if !seen[key] {
			seen[key] = true
			results = append(results, sr)
		}
	}

	// Limit to 8 results
	if len(results) > 8 {
		results = results[:8]
	}

	middleware.WriteJSON(w, http.StatusOK, results)
}

func (h *GeocodingHandler) searchPhoton(r *http.Request, query string) []searchResult {
	searchURL := fmt.Sprintf(
		"https://photon.komoot.io/api/?q=%s&limit=5&lat=6.2&lon=-75.6",
		url.QueryEscape(query),
	)

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, searchURL, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "CalculadoraPanelesSolares/1.0")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var photonResp photonResponse
	if err := json.Unmarshal(body, &photonResp); err != nil {
		return nil
	}

	results := make([]searchResult, 0, len(photonResp.Features))
	for _, f := range photonResp.Features {
		p := f.Properties
		parts := []string{}
		if p.Street != "" {
			street := p.Street
			if p.HouseNumber != "" {
				street += " #" + p.HouseNumber
			}
			parts = append(parts, street)
		} else if p.Name != "" {
			parts = append(parts, p.Name)
		}
		if p.City != "" {
			parts = append(parts, p.City)
		} else if p.County != "" {
			parts = append(parts, p.County)
		}
		if p.State != "" {
			parts = append(parts, p.State)
		}
		displayName := strings.Join(parts, ", ")
		if displayName == "" {
			continue
		}

		results = append(results, searchResult{
			DisplayName: displayName,
			Lat:         f.Geometry.Coordinates[1],
			Lon:         f.Geometry.Coordinates[0],
		})
	}
	return results
}

func (h *GeocodingHandler) searchNominatim(r *http.Request, query string) []searchResult {
	// Rate limiting for Nominatim
	h.mu.Lock()
	elapsed := time.Since(h.lastCall)
	if elapsed < time.Second {
		time.Sleep(time.Second - elapsed)
	}
	h.lastCall = time.Now()
	h.mu.Unlock()

	searchURL := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?q=%s&format=json&accept-language=es&countrycodes=co&limit=5&addressdetails=1",
		url.QueryEscape(query),
	)

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, searchURL, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "CalculadoraPanelesSolares/1.0")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var nomResults []nominatimSearchResult
	if err := json.Unmarshal(body, &nomResults); err != nil {
		return nil
	}

	results := make([]searchResult, 0, len(nomResults))
	for _, nr := range nomResults {
		lat, err1 := strconv.ParseFloat(nr.Lat, 64)
		lon, err2 := strconv.ParseFloat(nr.Lon, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		results = append(results, searchResult{
			DisplayName: nr.DisplayName,
			Lat:         lat,
			Lon:         lon,
		})
	}
	return results
}
