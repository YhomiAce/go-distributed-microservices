package main

import (
	"context"
	"logger/data"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)


type RPCServer struct {
	DB *mongo.Client
}

type RPCPayload struct {
	Data string
	Name string
}

func (r *RPCServer) LogInfo (payload RPCPayload, response *string) error {
	entry := data.LogEntry {
		Data: payload.Data,
		Name: payload.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := r.DB.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		return err
	}
	return nil
}