package repository

import (
	"context"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	colReq *mongo.Collection
	colRes *mongo.Collection
}

func NewRepository(colReq, colRes *mongo.Collection) *Repo {
	return &Repo{
		colReq: colReq,
		colRes: colRes}
}

func (r *Repo) GetAll() ([]models.Request, error) {
	cursor, err := r.colReq.Find(context.Background(), bson.D{}, options.Find())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var requests []models.Request

	for cursor.Next(context.Background()) {
		var request models.Request
		if err := cursor.Decode(&request); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (r *Repo) GetRequest(id string) (models.Request, error) {
	var request models.Request

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return request, err
	}

	filter := bson.M{"_id": objectID}

	err = r.colReq.FindOne(context.Background(), filter).Decode(&request)
	if err != nil {
		return request, err
	}

	return request, nil
}

//func (r *Repo) GetResponse(id int) (models.Response, error) {
//	return models.Response{}, nil
//}

func (r *Repo) SaveRequest(request models.Request) (primitive.ObjectID, error) {
	res, err := r.colReq.InsertOne(context.Background(), request)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	}

	return primitive.NilObjectID, nil
}

func (r *Repo) SaveResponse(response models.Response) error {
	_, err := r.colRes.InsertOne(context.Background(), response)
	if err != nil {
		return err
	}

	return nil
}
