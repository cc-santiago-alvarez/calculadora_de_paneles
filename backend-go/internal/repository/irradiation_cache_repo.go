package repository

import (
	"context"
	"time"

	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IrradiationCacheRepo struct {
	col *mongo.Collection
}

func NewIrradiationCacheRepo(db *mongo.Database) *IrradiationCacheRepo {
	return &IrradiationCacheRepo{col: db.Collection("irradiationcaches")}
}

func (r *IrradiationCacheRepo) FindByLocationKey(ctx context.Context, locationKey string) (*model.IrradiationCache, error) {
	var cache model.IrradiationCache
	err := r.col.FindOne(ctx, bson.M{
		"locationKey": locationKey,
		"expiresAt":   bson.M{"$gt": time.Now()},
	}).Decode(&cache)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &cache, nil
}

func (r *IrradiationCacheRepo) Upsert(ctx context.Context, locationKey, source string, data model.NormalizedGHI, ttlDays int) error {
	expiresAt := time.Now().AddDate(0, 0, ttlDays)
	opts := options.FindOneAndUpdate().SetUpsert(true)
	return r.col.FindOneAndUpdate(ctx, bson.M{"locationKey": locationKey}, bson.M{
		"$set": bson.M{
			"locationKey": locationKey,
			"source":      source,
			"fetchedAt":   time.Now(),
			"expiresAt":   expiresAt,
			"normalized":  data,
		},
	}, opts).Err()
}
