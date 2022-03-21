package db

import (
	"context"
	"fmt"
	"restapi-lesson/internal/user"
	"restapi-lesson/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	coolection *mongo.Collection
	logger     *logging.Logger
}

// Create implements user.Storage
func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.coolection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failde to create user due to error: %v", err)
	}

	d.logger.Debug("convert insertedID to objectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert objectid to hex. probably oid: %s", oid)
}

// Delete implements user.Storage
func (d *db) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindOne implements user.Storage
func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to objectid. hex: %s", id)
	}
	//mongo.getDatabase("test").getCollection("doc").find({})
	filter := bson.M{"_id": oid}
	result := d.coolection.FindOne(ctx, filter)
	if result.Err() != nil {
		//TODO 404
		return u, fmt.Errorf("failde to find one user by id: %s due to error: %v", id, err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user(id: %s) from db: %s due to error: %v", id, err)
	}
	return u, nil
}

// Update implements user.Storage
func (d *db) Update(ctx context.Context, user user.User) error {
	panic("unimplemented")
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {

	return &db{
		coolection: database.Collection(collection),
		Logger:     logger,
	}
}
