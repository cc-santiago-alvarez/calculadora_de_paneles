package model

import "go.mongodb.org/mongo-driver/v2/bson"

type ChargeControllerCatalog struct {
	ID                bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Manufacturer      string        `json:"manufacturer" bson:"manufacturer"`
	Model             string        `json:"model" bson:"model"`
	Type              string        `json:"type" bson:"type"` // "mppt" or "pwm"
	RatedPowerW       float64       `json:"ratedPowerW" bson:"ratedPowerW"`
	MaxPVVoltage      float64       `json:"maxPVVoltage" bson:"maxPVVoltage"`
	MaxPVCurrent      float64       `json:"maxPVCurrent" bson:"maxPVCurrent"`
	BatteryVoltages   []int         `json:"batteryVoltages" bson:"batteryVoltages"`
	MaxChargeCurrentA float64       `json:"maxChargeCurrentA" bson:"maxChargeCurrentA"`
	Efficiency        float64       `json:"efficiency" bson:"efficiency"`
	CostCOP           float64       `json:"costCOP" bson:"costCOP"`
	Warranty          int           `json:"warranty" bson:"warranty"`
	Weight            float64       `json:"weight,omitempty" bson:"weight,omitempty"`
	IsActive          bool          `json:"isActive" bson:"isActive"`
	Source            string        `json:"source,omitempty" bson:"source,omitempty"`
}
