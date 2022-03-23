package db

import (
	"context"
	"errors"
	"fmt"
	"restapi-lesson/internal/apperror"
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
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", id)
	}
	filter := bson.M{"_id": objectID}
	result, err := d.coolection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failde to execute query. %v", err)
	}

	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}
	d.logger.Trace("Deleted %d documents", result.DeletedCount)

	return nil
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
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.ErrNotFound
		}
		return u, fmt.Errorf("failde to find one user by id: %s due to error: %v", id, err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user(id: %s) from DB due to error: %v", id, err)
	}
	return u, nil
}

// Update implements user.Storage
func (d *db) Update(ctx context.Context, user user.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", user.ID)
	}
	filter := bson.M{"_id": objectID}
	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user. error: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user bytes. error: %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.coolection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query. error: %v", err)
	}

	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}

	d.logger.Tracef("Matched %d documents and modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (d *db) FindAll(ctx context.Context) (u []user.User, err error) {
	cursor, err := d.coolection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return u, fmt.Errorf("failde to find all users due to error: %v", err)
	}

	err = cursor.All(ctx, &u)
	if err != nil {
		return u, fmt.Errorf("failed to read all documents from cursor. error: %v", err)
	}

	return u, nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {

	return &db{
		coolection: database.Collection(collection),
		logger:     logger,
	}
}
