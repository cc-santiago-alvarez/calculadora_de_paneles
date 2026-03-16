package repository

import (
	"context"

	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CatalogRepo struct {
	panels    *mongo.Collection
	inverters *mongo.Collection
}

func NewCatalogRepo(db *mongo.Database) *CatalogRepo {
	return &CatalogRepo{
		panels:    db.Collection("panelcatalogs"),
		inverters: db.Collection("invertercatalogs"),
	}
}

func (r *CatalogRepo) FindPanels(ctx context.Context, filter bson.M) ([]model.PanelCatalog, error) {
	if filter == nil {
		filter = bson.M{}
	}
	filter["isActive"] = true
	opts := options.Find().SetSort(bson.D{{Key: "powerWp", Value: -1}})
	cursor, err := r.panels.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var panels []model.PanelCatalog
	if err := cursor.All(ctx, &panels); err != nil {
		return nil, err
	}
	if panels == nil {
		panels = []model.PanelCatalog{}
	}
	return panels, nil
}

func (r *CatalogRepo) FindPanelByID(ctx context.Context, id bson.ObjectID) (*model.PanelCatalog, error) {
	var panel model.PanelCatalog
	err := r.panels.FindOne(ctx, bson.M{"_id": id}).Decode(&panel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &panel, nil
}

func (r *CatalogRepo) FindInverters(ctx context.Context, filter bson.M) ([]model.InverterCatalog, error) {
	if filter == nil {
		filter = bson.M{}
	}
	filter["isActive"] = true
	opts := options.Find().SetSort(bson.D{{Key: "ratedPowerKw", Value: 1}})
	cursor, err := r.inverters.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var inverters []model.InverterCatalog
	if err := cursor.All(ctx, &inverters); err != nil {
		return nil, err
	}
	if inverters == nil {
		inverters = []model.InverterCatalog{}
	}
	return inverters, nil
}

func (r *CatalogRepo) FindInverterByID(ctx context.Context, id bson.ObjectID) (*model.InverterCatalog, error) {
	var inverter model.InverterCatalog
	err := r.inverters.FindOne(ctx, bson.M{"_id": id}).Decode(&inverter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &inverter, nil
}

// UpsertPanels upserts panels by manufacturer+model as unique key.
func (r *CatalogRepo) UpsertPanels(ctx context.Context, panels []model.PanelCatalog) (int, error) {
	upserted := 0
	for _, p := range panels {
		filter := bson.M{"manufacturer": p.Manufacturer, "model": p.Model}
		update := bson.M{"$set": p}
		opts := options.UpdateOne().SetUpsert(true)
		result, err := r.panels.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return upserted, err
		}
		if result.UpsertedCount > 0 || result.ModifiedCount > 0 {
			upserted++
		}
	}
	return upserted, nil
}

// UpsertInverters upserts inverters by manufacturer+model as unique key.
func (r *CatalogRepo) UpsertInverters(ctx context.Context, inverters []model.InverterCatalog) (int, error) {
	upserted := 0
	for _, inv := range inverters {
		filter := bson.M{"manufacturer": inv.Manufacturer, "model": inv.Model}
		update := bson.M{"$set": inv}
		opts := options.UpdateOne().SetUpsert(true)
		result, err := r.inverters.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return upserted, err
		}
		if result.UpsertedCount > 0 || result.ModifiedCount > 0 {
			upserted++
		}
	}
	return upserted, nil
}

// SeedPanels inserts default panels if collection is empty.
func (r *CatalogRepo) SeedPanels(ctx context.Context, panels []model.PanelCatalog) error {
	count, err := r.panels.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	docs := make([]interface{}, len(panels))
	for i, p := range panels {
		docs[i] = p
	}
	_, err = r.panels.InsertMany(ctx, docs)
	return err
}

// SeedInverters inserts default inverters if collection is empty.
func (r *CatalogRepo) SeedInverters(ctx context.Context, inverters []model.InverterCatalog) error {
	count, err := r.inverters.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	docs := make([]interface{}, len(inverters))
	for i, inv := range inverters {
		docs[i] = inv
	}
	_, err = r.inverters.InsertMany(ctx, docs)
	return err
}
