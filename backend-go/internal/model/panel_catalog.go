package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Dimensions struct {
	Length float64 `json:"length" bson:"length"`
	Width  float64 `json:"width" bson:"width"`
}

type PanelCatalog struct {
	ID             bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Manufacturer   string        `json:"manufacturer" bson:"manufacturer"`
	Model          string        `json:"model" bson:"model"`
	Type           string        `json:"type" bson:"type"`
	PowerWp        float64       `json:"powerWp" bson:"powerWp"`
	Efficiency     float64       `json:"efficiency" bson:"efficiency"`
	Area           float64       `json:"area" bson:"area"`
	Voc            float64       `json:"voc" bson:"voc"`
	Isc            float64       `json:"isc" bson:"isc"`
	Vmp            float64       `json:"vmp" bson:"vmp"`
	Imp            float64       `json:"imp" bson:"imp"`
	TempCoeffPmax  float64       `json:"tempCoeffPmax" bson:"tempCoeffPmax"`
	TempCoeffVoc   float64       `json:"tempCoeffVoc" bson:"tempCoeffVoc"`
	NOCT           float64       `json:"NOCT" bson:"NOCT"`
	Weight         float64       `json:"weight,omitempty" bson:"weight,omitempty"`
	PanelDimensions Dimensions   `json:"dimensions,omitempty" bson:"dimensions,omitempty"`
	Warranty       int           `json:"warranty" bson:"warranty"`
	CostCOP        float64       `json:"costCOP" bson:"costCOP"`
	IsActive       bool          `json:"isActive" bson:"isActive"`
	Source         string        `json:"source,omitempty" bson:"source,omitempty"`
	CecID          string        `json:"cecId,omitempty" bson:"cecId,omitempty"`
}
