package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Location struct {
	Latitude    float64 `json:"latitude" bson:"latitude"`
	Longitude   float64 `json:"longitude" bson:"longitude"`
	Altitude    float64 `json:"altitude" bson:"altitude"`
	ClimateZone string  `json:"climateZone" bson:"climateZone"`
	Department  string  `json:"department" bson:"department"`
	City        string  `json:"city" bson:"city"`
}

type Consumption struct {
	Monthly        [12]float64 `json:"monthly" bson:"monthly"`
	TariffPerKwh   float64     `json:"tariffPerKwh" bson:"tariffPerKwh"`
	Estrato        int         `json:"estrato" bson:"estrato"`
	ConnectionType string      `json:"connectionType" bson:"connectionType"`
}

type ShadingProfile struct {
	HasShading  bool        `json:"hasShading" bson:"hasShading"`
	MonthlyLoss [12]float64 `json:"monthlyLoss" bson:"monthlyLoss"`
}

type Roof struct {
	Area             float64        `json:"area" bson:"area"`
	Azimuth          float64        `json:"azimuth" bson:"azimuth"`
	Tilt             float64        `json:"tilt" bson:"tilt"`
	UsablePercentage float64        `json:"usablePercentage" bson:"usablePercentage"`
	ShadingProfile   ShadingProfile `json:"shadingProfile" bson:"shadingProfile"`
}

type PanelOverride struct {
	Watts float64 `json:"watts,omitempty" bson:"watts,omitempty"`
	Area  float64 `json:"area,omitempty" bson:"area,omitempty"`
}

type Equipment struct {
	PanelID       bson.ObjectID  `json:"panelId" bson:"panelId"`
	InverterID    bson.ObjectID  `json:"inverterId" bson:"inverterId"`
	PanelOverride *PanelOverride `json:"panelOverride,omitempty" bson:"panelOverride,omitempty"`
}

type Project struct {
	ID                 bson.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name               string          `json:"name" bson:"name"`
	Location           Location        `json:"location" bson:"location"`
	Consumption        Consumption     `json:"consumption" bson:"consumption"`
	Roof               Roof            `json:"roof" bson:"roof"`
	SystemType         string          `json:"systemType" bson:"systemType"`
	CoveragePercentage float64         `json:"coveragePercentage" bson:"coveragePercentage"`
	Equipment          Equipment       `json:"equipment" bson:"equipment"`
	Scenarios          []bson.ObjectID `json:"scenarios" bson:"scenarios"`
	CreatedAt          time.Time       `json:"createdAt" bson:"createdAt"`
	UpdatedAt          time.Time       `json:"updatedAt" bson:"updatedAt"`
}
