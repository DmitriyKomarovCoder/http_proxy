package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewMongoConnect(ctx context.Context, url, databaseName, colRequest, colResponse string) (*mongo.Collection, *mongo.Collection, *mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(url)
	clientOptions.SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, nil, err
	}

	colReq := client.Database(databaseName).Collection(colRequest)
	colRes := client.Database(databaseName).Collection(colResponse)

	_, err = colRes.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.M{"_id": 1},
		},
	)

	if err != nil {
		return nil, nil, nil, err
	}

	_, err = colReq.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.M{"_id": 1},
		},
	)

	if err != nil {
		return nil, nil, nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	return colReq, colRes, client, nil
}
