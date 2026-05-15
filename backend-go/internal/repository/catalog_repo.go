package repository

import (
	"context"
	"errors"

	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ErrDuplicatePanel is returned when creating a panel that already exists
// (same manufacturer + model).
var ErrDuplicatePanel = errors.New("ya existe un panel con ese fabricante y modelo")

// ErrDuplicateInverter is returned when creating an inverter that already exists
// (same manufacturer + model).
var ErrDuplicateInverter = errors.New("ya existe un inversor con ese fabricante y modelo")

type CatalogRepo struct {
	panels            *mongo.Collection
	inverters         *mongo.Collection
	chargeControllers *mongo.Collection
}

func NewCatalogRepo(db *mongo.Database) *CatalogRepo {
	return &CatalogRepo{
		panels:            db.Collection("panelcatalogs"),
		inverters:         db.Collection("invertercatalogs"),
		chargeControllers: db.Collection("chargecontrollercatalogs"),
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

// CreatePanel inserts a new panel into the catalog. Returns ErrDuplicatePanel
// if a panel with the same manufacturer+model already exists.
func (r *CatalogRepo) CreatePanel(ctx context.Context, panel model.PanelCatalog) (*model.PanelCatalog, error) {
	existing := r.panels.FindOne(ctx, bson.M{"manufacturer": panel.Manufacturer, "model": panel.Model})
	if err := existing.Err(); err == nil {
		return nil, ErrDuplicatePanel
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	panel.ID = bson.NilObjectID
	panel.IsActive = true
	result, err := r.panels.InsertOne(ctx, panel)
	if err != nil {
		return nil, err
	}
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		panel.ID = oid
	}
	return &panel, nil
}

// UpdatePanel updates an existing panel by ID and returns the updated document.
func (r *CatalogRepo) UpdatePanel(ctx context.Context, id bson.ObjectID, panel model.PanelCatalog) (*model.PanelCatalog, error) {
	update := bson.M{"$set": bson.M{
		"manufacturer":  panel.Manufacturer,
		"model":         panel.Model,
		"type":          panel.Type,
		"powerWp":       panel.PowerWp,
		"efficiency":    panel.Efficiency,
		"area":          panel.Area,
		"voc":           panel.Voc,
		"isc":           panel.Isc,
		"vmp":           panel.Vmp,
		"imp":           panel.Imp,
		"tempCoeffPmax": panel.TempCoeffPmax,
		"tempCoeffVoc":  panel.TempCoeffVoc,
		"NOCT":          panel.NOCT,
		"weight":        panel.Weight,
		"dimensions":    panel.PanelDimensions,
		"warranty":      panel.Warranty,
		"costCOP":       panel.CostCOP,
		"format":        panel.Format,
	}}
	result, err := r.panels.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, nil
	}
	return r.FindPanelByID(ctx, id)
}

// SoftDeletePanel marks a panel as inactive instead of removing it.
func (r *CatalogRepo) SoftDeletePanel(ctx context.Context, id bson.ObjectID) (bool, error) {
	result, err := r.panels.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"isActive": false}})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, nil
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

// CreateInverter inserts a new inverter into the catalog. Returns
// ErrDuplicateInverter if an inverter with the same manufacturer+model exists.
func (r *CatalogRepo) CreateInverter(ctx context.Context, inverter model.InverterCatalog) (*model.InverterCatalog, error) {
	existing := r.inverters.FindOne(ctx, bson.M{"manufacturer": inverter.Manufacturer, "model": inverter.Model})
	if err := existing.Err(); err == nil {
		return nil, ErrDuplicateInverter
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	inverter.ID = bson.NilObjectID
	inverter.IsActive = true
	result, err := r.inverters.InsertOne(ctx, inverter)
	if err != nil {
		return nil, err
	}
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		inverter.ID = oid
	}
	return &inverter, nil
}

// UpdateInverter updates an existing inverter by ID and returns the updated document.
func (r *CatalogRepo) UpdateInverter(ctx context.Context, id bson.ObjectID, inverter model.InverterCatalog) (*model.InverterCatalog, error) {
	update := bson.M{"$set": bson.M{
		"manufacturer":    inverter.Manufacturer,
		"model":           inverter.Model,
		"type":            inverter.Type,
		"ratedPowerKw":    inverter.RatedPowerKw,
		"maxDCPowerKw":    inverter.MaxDCPowerKw,
		"efficiency":      inverter.Efficiency,
		"mpptCount":       inverter.MPPTCount,
		"mpptVoltageMin":  inverter.MPPTVoltageMin,
		"mpptVoltageMax":  inverter.MPPTVoltageMax,
		"maxInputVoltage": inverter.MaxInputVoltage,
		"maxInputCurrent": inverter.MaxInputCurrent,
		"outputVoltage":   inverter.OutputVoltage,
		"outputPhases":    inverter.OutputPhases,
		"hasBatteryPort":  inverter.HasBatteryPort,
		"weight":          inverter.Weight,
		"warranty":        inverter.Warranty,
		"costCOP":         inverter.CostCOP,
	}}
	result, err := r.inverters.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, nil
	}
	return r.FindInverterByID(ctx, id)
}

// SoftDeleteInverter marks an inverter as inactive instead of removing it.
func (r *CatalogRepo) SoftDeleteInverter(ctx context.Context, id bson.ObjectID) (bool, error) {
	result, err := r.inverters.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"isActive": false}})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, nil
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

// FindChargeControllers returns charge controllers matching the filter.
func (r *CatalogRepo) FindChargeControllers(ctx context.Context, filter bson.M) ([]model.ChargeControllerCatalog, error) {
	if filter == nil {
		filter = bson.M{}
	}
	filter["isActive"] = true
	opts := options.Find().SetSort(bson.D{{Key: "maxChargeCurrentA", Value: 1}})
	cursor, err := r.chargeControllers.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var controllers []model.ChargeControllerCatalog
	if err := cursor.All(ctx, &controllers); err != nil {
		return nil, err
	}
	if controllers == nil {
		controllers = []model.ChargeControllerCatalog{}
	}
	return controllers, nil
}

// FindChargeControllerByID returns a single charge controller by ID.
func (r *CatalogRepo) FindChargeControllerByID(ctx context.Context, id bson.ObjectID) (*model.ChargeControllerCatalog, error) {
	var cc model.ChargeControllerCatalog
	err := r.chargeControllers.FindOne(ctx, bson.M{"_id": id}).Decode(&cc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &cc, nil
}

// UpsertChargeControllers upserts charge controllers by manufacturer+model as unique key.
func (r *CatalogRepo) UpsertChargeControllers(ctx context.Context, controllers []model.ChargeControllerCatalog) (int, error) {
	upserted := 0
	for _, cc := range controllers {
		filter := bson.M{"manufacturer": cc.Manufacturer, "model": cc.Model}
		update := bson.M{"$set": cc}
		opts := options.UpdateOne().SetUpsert(true)
		result, err := r.chargeControllers.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return upserted, err
		}
		if result.UpsertedCount > 0 || result.ModifiedCount > 0 {
			upserted++
		}
	}
	return upserted, nil
}
