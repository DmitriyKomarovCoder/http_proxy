package api

import (
	"github.com/DmitriyKomarovCoder/http_proxy/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Usecase interface {
	AllRequest() ([]models.Request, error)
	GetRequest(id string) (models.Request, error)
	Repeat(id string) (*http.Request, error)
	Scan(id string) (bool, error)
	SaveRequest(request *http.Request, bodyBytes []byte) (primitive.ObjectID, error)
	SaveResponse(id primitive.ObjectID, response http.Response, body []byte) error
}

type Repository interface {
	GetAll() ([]models.Request, error)
	GetRequest(id string) (models.Request, error)
	SaveRequest(request models.Request) (primitive.ObjectID, error)
	SaveResponse(response models.Response) error
}
