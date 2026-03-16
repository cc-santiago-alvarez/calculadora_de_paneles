package repository

import (
	"context"
	"time"

	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProjectRepo struct {
	col *mongo.Collection
}

func NewProjectRepo(db *mongo.Database) *ProjectRepo {
	return &ProjectRepo{col: db.Collection("projects")}
}

func (r *ProjectRepo) FindAll(ctx context.Context) ([]model.Project, error) {
	opts := options.Find().SetSort(bson.D{{Key: "updatedAt", Value: -1}})
	cursor, err := r.col.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	var projects []model.Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}
	if projects == nil {
		projects = []model.Project{}
	}
	return projects, nil
}

func (r *ProjectRepo) FindByID(ctx context.Context, id bson.ObjectID) (*model.Project, error) {
	var project model.Project
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) Create(ctx context.Context, project *model.Project) error {
	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now
	if project.Scenarios == nil {
		project.Scenarios = []bson.ObjectID{}
	}
	result, err := r.col.InsertOne(ctx, project)
	if err != nil {
		return err
	}
	project.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

func (r *ProjectRepo) Update(ctx context.Context, id bson.ObjectID, update bson.M) (*model.Project, error) {
	update["updatedAt"] = time.Now()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var project model.Project
	err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": update}, opts).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) Delete(ctx context.Context, id bson.ObjectID) error {
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *ProjectRepo) PushScenario(ctx context.Context, projectID, scenarioID bson.ObjectID) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": projectID}, bson.M{
		"$push": bson.M{"scenarios": scenarioID},
		"$set":  bson.M{"updatedAt": time.Now()},
	})
	return err
}
