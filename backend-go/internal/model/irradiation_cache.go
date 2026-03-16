package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NormalizedGHI struct {
	MonthlyGHI []float64 `json:"monthlyGHI" bson:"monthlyGHI"`
	AnnualGHI  float64   `json:"annualGHI" bson:"annualGHI"`
	Elevation  float64   `json:"elevation" bson:"elevation"`
}

type IrradiationCache struct {
	ID          bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	LocationKey string        `json:"locationKey" bson:"locationKey"`
	Source      string        `json:"source" bson:"source"`
	FetchedAt   time.Time     `json:"fetchedAt" bson:"fetchedAt"`
	ExpiresAt   time.Time     `json:"expiresAt" bson:"expiresAt"`
	Normalized  NormalizedGHI `json:"normalized" bson:"normalized"`
}
