package handler

import (
	"net/http"

	"github.com/dev13/calculadora-paneles-backend/internal/middleware"
	"github.com/dev13/calculadora-paneles-backend/internal/service"
)

type IrradiationHandler struct {
	irradiationService *service.IrradiationService
}

func NewIrradiationHandler(irradiationService *service.IrradiationService) *IrradiationHandler {
	return &IrradiationHandler{irradiationService: irradiationService}
}

func (h *IrradiationHandler) FetchIrradiation(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
		Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
		Tilt      float64 `json:"tilt"`
		Azimuth   float64 `json:"azimuth"`
	}

	if !middleware.ValidateAndDecode(w, r, &input) {
		return
	}

	if input.Tilt == 0 {
		input.Tilt = 10
	}

	result, err := h.irradiationService.GetIrradiation(r.Context(), input.Latitude, input.Longitude, input.Tilt, input.Azimuth, 0.2)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusOK, result)
}
