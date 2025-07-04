package db

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.elastic.co/apm/module/apmmongo/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mongoptions "go.mongodb.org/mongo-driver/mongo/options"
)

type (
	MongoOptions struct {
		URL        string
		DisableAPM bool
	}

	Mongo struct {
		mongo.Client
		dbName string
	}
)

func (m *Mongo) Collection(name string) *mongo.Collection {
	return m.Database(m.dbName).Collection(name)
}

func OpenMongo(ctx context.Context, uri string) (*Mongo, error) {
	return OpenMongoOpts(ctx, MongoOptions{
		URL:        uri,
		DisableAPM: false,
	})
}

func OpenMongoOpts(ctx context.Context, opts MongoOptions) (*Mongo, error) {
	mongoOptions := mongoptions.Client().ApplyURI(opts.URL)
	if !opts.DisableAPM {
		mongoOptions = mongoOptions.SetMonitor(apmmongo.CommandMonitor())
	}

	cli, err := mongo.Connect(ctx, mongoOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	// Parse the URI to get the database name
	parsedURI, err := url.Parse(opts.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the mongo URI: %w", err)
	}

	dbName := strings.TrimPrefix(parsedURI.Path, "/")
	if dbName == "" {
		return nil, fmt.Errorf("failed to get the database name from the mongo URI")
	}
	return &Mongo{Client: *cli, dbName: dbName}, nil
}

func MongoParseObjectID(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("MongoParseObjectID: %w", err)
	}
	return objectID, nil
}

func MongoParseTsFromObjectID(id string) (time.Time, error) {
	objectID, err := MongoParseObjectID(id)
	if err != nil {
		return time.Time{}, fmt.Errorf("MongoParseTsFromObjectID: %w", err)
	}
	return objectID.Timestamp().UTC(), nil
}
