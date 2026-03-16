package repository

import (
	"context"
	"time"

	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ScenarioRepo struct {
	col *mongo.Collection
}

func NewScenarioRepo(db *mongo.Database) *ScenarioRepo {
	return &ScenarioRepo{col: db.Collection("scenarios")}
}

func (r *ScenarioRepo) FindByProjectID(ctx context.Context, projectID bson.ObjectID) ([]model.Scenario, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.col.Find(ctx, bson.M{"projectId": projectID}, opts)
	if err != nil {
		return nil, err
	}
	var scenarios []model.Scenario
	if err := cursor.All(ctx, &scenarios); err != nil {
		return nil, err
	}
	if scenarios == nil {
		scenarios = []model.Scenario{}
	}
	return scenarios, nil
}

func (r *ScenarioRepo) FindByID(ctx context.Context, id bson.ObjectID) (*model.Scenario, error) {
	var scenario model.Scenario
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&scenario)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &scenario, nil
}

func (r *ScenarioRepo) Create(ctx context.Context, scenario *model.Scenario) error {
	now := time.Now()
	scenario.CreatedAt = now
	scenario.UpdatedAt = now
	result, err := r.col.InsertOne(ctx, scenario)
	if err != nil {
		return err
	}
	scenario.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

func (r *ScenarioRepo) DeleteByProjectID(ctx context.Context, projectID bson.ObjectID) error {
	_, err := r.col.DeleteMany(ctx, bson.M{"projectId": projectID})
	return err
}
