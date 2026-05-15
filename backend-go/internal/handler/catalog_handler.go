package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/dev13/calculadora-paneles-backend/internal/middleware"
	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"github.com/dev13/calculadora-paneles-backend/internal/repository"
	"github.com/dev13/calculadora-paneles-backend/internal/service"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CatalogHandler struct {
	catalogRepo *repository.CatalogRepo
	cecService  *service.CECCatalogService
}

func NewCatalogHandler(catalogRepo *repository.CatalogRepo, cecService *service.CECCatalogService) *CatalogHandler {
	return &CatalogHandler{
		catalogRepo: catalogRepo,
		cecService:  cecService,
	}
}

func (h *CatalogHandler) GetPanels(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{}
	q := r.URL.Query()
	if t := q.Get("type"); t != "" {
		filter["type"] = t
	}
	if minP := q.Get("minPower"); minP != "" {
		if v, err := strconv.ParseFloat(minP, 64); err == nil {
			if existing, ok := filter["powerWp"].(bson.M); ok {
				existing["$gte"] = v
			} else {
				filter["powerWp"] = bson.M{"$gte": v}
			}
		}
	}
	if maxP := q.Get("maxPower"); maxP != "" {
		if v, err := strconv.ParseFloat(maxP, 64); err == nil {
			if existing, ok := filter["powerWp"].(bson.M); ok {
				existing["$lte"] = v
			} else {
				filter["powerWp"] = bson.M{"$lte": v}
			}
		}
	}
	if mfr := q.Get("manufacturer"); mfr != "" {
		filter["manufacturer"] = bson.M{"$regex": mfr, "$options": "i"}
	}

	panels, err := h.catalogRepo.FindPanels(r.Context(), filter)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusOK, panels)
}

func (h *CatalogHandler) GetInverters(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{}
	q := r.URL.Query()
	if t := q.Get("type"); t != "" {
		filter["type"] = t
	}
	if minP := q.Get("minPower"); minP != "" {
		if v, err := strconv.ParseFloat(minP, 64); err == nil {
			if existing, ok := filter["ratedPowerKw"].(bson.M); ok {
				existing["$gte"] = v
			} else {
				filter["ratedPowerKw"] = bson.M{"$gte": v}
			}
		}
	}
	if maxP := q.Get("maxPower"); maxP != "" {
		if v, err := strconv.ParseFloat(maxP, 64); err == nil {
			if existing, ok := filter["ratedPowerKw"].(bson.M); ok {
				existing["$lte"] = v
			} else {
				filter["ratedPowerKw"] = bson.M{"$lte": v}
			}
		}
	}
	if mfr := q.Get("manufacturer"); mfr != "" {
		filter["manufacturer"] = bson.M{"$regex": mfr, "$options": "i"}
	}
	if hasBat := q.Get("hasBattery"); hasBat == "true" {
		filter["hasBatteryPort"] = true
	} else if hasBat == "false" {
		filter["hasBatteryPort"] = false
	}

	inverters, err := h.catalogRepo.FindInverters(r.Context(), filter)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusOK, inverters)
}

func (h *CatalogHandler) GetChargeControllers(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{}
	q := r.URL.Query()
	if t := q.Get("type"); t != "" {
		filter["type"] = t
	}
	if minC := q.Get("minCurrent"); minC != "" {
		if v, err := strconv.ParseFloat(minC, 64); err == nil {
			if existing, ok := filter["maxChargeCurrentA"].(bson.M); ok {
				existing["$gte"] = v
			} else {
				filter["maxChargeCurrentA"] = bson.M{"$gte": v}
			}
		}
	}
	if maxC := q.Get("maxCurrent"); maxC != "" {
		if v, err := strconv.ParseFloat(maxC, 64); err == nil {
			if existing, ok := filter["maxChargeCurrentA"].(bson.M); ok {
				existing["$lte"] = v
			} else {
				filter["maxChargeCurrentA"] = bson.M{"$lte": v}
			}
		}
	}
	if mfr := q.Get("manufacturer"); mfr != "" {
		filter["manufacturer"] = bson.M{"$regex": mfr, "$options": "i"}
	}
	if bv := q.Get("batteryVoltage"); bv != "" {
		if v, err := strconv.Atoi(bv); err == nil {
			filter["batteryVoltages"] = v
		}
	}

	controllers, err := h.catalogRepo.FindChargeControllers(r.Context(), filter)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusOK, controllers)
}

// GetPanelByID returns a single panel from the catalog by its ID.
func (h *CatalogHandler) GetPanelByID(w http.ResponseWriter, r *http.Request) {
	id, err := bson.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, "ID inválido"))
		return
	}
	panel, err := h.catalogRepo.FindPanelByID(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	if panel == nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusNotFound, "Panel no encontrado"))
		return
	}
	middleware.WriteJSON(w, http.StatusOK, panel)
}

// CreatePanel adds a new panel to the catalog.
func (h *CatalogHandler) CreatePanel(w http.ResponseWriter, r *http.Request) {
	var panel model.PanelCatalog
	if err := middleware.ReadJSON(r, &panel); err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, "JSON inválido"))
		return
	}
	if msg := validatePanel(panel); msg != "" {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, msg))
		return
	}

	created, err := h.catalogRepo.CreatePanel(r.Context(), panel)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicatePanel) {
			middleware.WriteError(w, middleware.NewAppError(http.StatusConflict, err.Error()))
			return
		}
		middleware.WriteError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusCreated, created)
}

// UpdatePanel updates an existing panel in the catalog.
func (h *CatalogHandler) UpdatePanel(w http.ResponseWriter, r *http.Request) {
	id, err := bson.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, "ID inválido"))
		return
	}
	var panel model.PanelCatalog
	if err := middleware.ReadJSON(r, &panel); err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, "JSON inválido"))
		return
	}
	if msg := validatePanel(panel); msg != "" {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, msg))
		return
	}

	updated, err := h.catalogRepo.UpdatePanel(r.Context(), id, panel)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	if updated == nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusNotFound, "Panel no encontrado"))
		return
	}
	middleware.WriteJSON(w, http.StatusOK, updated)
}

// DeletePanel performs a soft delete (isActive=false) of a panel.
func (h *CatalogHandler) DeletePanel(w http.ResponseWriter, r *http.Request) {
	id, err := bson.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, "ID inválido"))
		return
	}
	found, err := h.catalogRepo.SoftDeletePanel(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	if !found {
		middleware.WriteError(w, middleware.NewAppError(http.StatusNotFound, "Panel no encontrado"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// validatePanel returns an error message if required fields are missing.
func validatePanel(p model.PanelCatalog) string {
	if strings.TrimSpace(p.Manufacturer) == "" {
		return "El fabricante es obligatorio"
	}
	if strings.TrimSpace(p.Model) == "" {
		return "El modelo es obligatorio"
	}
	if p.PowerWp <= 0 {
		return "La potencia (powerWp) debe ser mayor a 0"
	}
	return ""
}

// GetInverterByID returns a single inverter from the catalog by its ID.
func (h *CatalogHandler) GetInverterByID(w http.ResponseWriter, r *http.Request) {
	id, err := bson.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, "ID inválido"))
		return
	}
	inverter, err := h.catalogRepo.FindInverterByID(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	if inverter == nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusNotFound, "Inversor no encontrado"))
		return
	}
	middleware.WriteJSON(w, http.StatusOK, inverter)
}

// CreateInverter adds a new inverter to the catalog.
func (h *CatalogHandler) CreateInverter(w http.ResponseWriter, r *http.Request) {
	var inverter model.InverterCatalog
	if err := middleware.ReadJSON(r, &inverter); err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, "JSON inválido"))
		return
	}
	if msg := validateInverter(inverter); msg != "" {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, msg))
		return
	}

	created, err := h.catalogRepo.CreateInverter(r.Context(), inverter)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateInverter) {
			middleware.WriteError(w, middleware.NewAppError(http.StatusConflict, err.Error()))
			return
		}
		middleware.WriteError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusCreated, created)
}

// UpdateInverter updates an existing inverter in the catalog.
func (h *CatalogHandler) UpdateInverter(w http.ResponseWriter, r *http.Request) {
	id, err := bson.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, "ID inválido"))
		return
	}
	var inverter model.InverterCatalog
	if err := middleware.ReadJSON(r, &inverter); err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, "JSON inválido"))
		return
	}
	if msg := validateInverter(inverter); msg != "" {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, msg))
		return
	}

	updated, err := h.catalogRepo.UpdateInverter(r.Context(), id, inverter)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	if updated == nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusNotFound, "Inversor no encontrado"))
		return
	}
	middleware.WriteJSON(w, http.StatusOK, updated)
}

// DeleteInverter performs a soft delete (isActive=false) of an inverter.
func (h *CatalogHandler) DeleteInverter(w http.ResponseWriter, r *http.Request) {
	id, err := bson.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(http.StatusBadRequest, "ID inválido"))
		return
	}
	found, err := h.catalogRepo.SoftDeleteInverter(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	if !found {
		middleware.WriteError(w, middleware.NewAppError(http.StatusNotFound, "Inversor no encontrado"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// validateInverter returns an error message if required fields are missing.
func validateInverter(inv model.InverterCatalog) string {
	if strings.TrimSpace(inv.Manufacturer) == "" {
		return "El fabricante es obligatorio"
	}
	if strings.TrimSpace(inv.Model) == "" {
		return "El modelo es obligatorio"
	}
	if inv.RatedPowerKw <= 0 {
		return "La potencia (ratedPowerKw) debe ser mayor a 0"
	}
	return ""
}

func (h *CatalogHandler) SyncCEC(w http.ResponseWriter, r *http.Request) {
	result, err := h.cecService.SyncAll(r.Context())
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusOK, result)
}
