package model

import "go.mongodb.org/mongo-driver/v2/bson"

type InverterCatalog struct {
	ID              bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Manufacturer    string        `json:"manufacturer" bson:"manufacturer"`
	Model           string        `json:"model" bson:"model"`
	Type            string        `json:"type" bson:"type"`
	RatedPowerKw    float64       `json:"ratedPowerKw" bson:"ratedPowerKw"`
	MaxDCPowerKw    float64       `json:"maxDCPowerKw" bson:"maxDCPowerKw"`
	Efficiency      float64       `json:"efficiency" bson:"efficiency"`
	MPPTCount       int           `json:"mpptCount" bson:"mpptCount"`
	MPPTVoltageMin  float64       `json:"mpptVoltageMin" bson:"mpptVoltageMin"`
	MPPTVoltageMax  float64       `json:"mpptVoltageMax" bson:"mpptVoltageMax"`
	MaxInputVoltage float64       `json:"maxInputVoltage" bson:"maxInputVoltage"`
	MaxInputCurrent float64       `json:"maxInputCurrent" bson:"maxInputCurrent"`
	OutputVoltage   float64       `json:"outputVoltage" bson:"outputVoltage"`
	OutputPhases    int           `json:"outputPhases" bson:"outputPhases"`
	HasBatteryPort  bool          `json:"hasBatteryPort" bson:"hasBatteryPort"`
	Weight          float64       `json:"weight,omitempty" bson:"weight,omitempty"`
	Warranty        int           `json:"warranty" bson:"warranty"`
	CostCOP         float64       `json:"costCOP" bson:"costCOP"`
	IsActive        bool          `json:"isActive" bson:"isActive"`
	Source          string        `json:"source,omitempty" bson:"source,omitempty"`
	CecID           string        `json:"cecId,omitempty" bson:"cecId,omitempty"`
}
