package handler

import (
	"net/http"
	"strconv"

	"github.com/dev13/calculadora-paneles-backend/internal/middleware"
	"github.com/dev13/calculadora-paneles-backend/internal/repository"
	"github.com/dev13/calculadora-paneles-backend/internal/service"
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
	}

	inverters, err := h.catalogRepo.FindInverters(r.Context(), filter)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusOK, inverters)
}

func (h *CatalogHandler) SyncCEC(w http.ResponseWriter, r *http.Request) {
	result, err := h.cecService.SyncAll(r.Context())
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusOK, result)
}
