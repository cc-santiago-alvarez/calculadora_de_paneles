package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func ConnectMongoDB(ctx context.Context, uri string) (*MongoDB, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database("calculadora_paneles")
	log.Println("✓ MongoDB connected")

	return &MongoDB{Client: client, Database: db}, nil
}

func (m *MongoDB) EnsureIndexes(ctx context.Context) error {
	// IrradiationCache: unique locationKey index
	cacheCol := m.Database.Collection("irradiationcaches")
	_, err := cacheCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "locationKey", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create irradiation cache indexes: %w", err)
	}

	// PanelCatalog: unique manufacturer+model
	panelCol := m.Database.Collection("panelcatalogs")
	_, err = panelCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "manufacturer", Value: 1}, {Key: "model", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create panel catalog index: %w", err)
	}

	// InverterCatalog: unique manufacturer+model
	inverterCol := m.Database.Collection("invertercatalogs")
	_, err = inverterCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "manufacturer", Value: 1}, {Key: "model", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create inverter catalog index: %w", err)
	}

	return nil
}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
