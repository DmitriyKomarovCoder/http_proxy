package api

import "github.com/DmitriyKomarovCoder/http_proxy/internal/models"

type Usecase interface {
	AllRequest() ([]models.Request, error)
	GetRequest(id int) (models.Request, error)
	Repeat(id int) (models.Response, error)
	Scan(id int) error
}

type Repository interface {
	GetAll() ([]models.Request, error)
	GetRequest(id int) (models.Request, error)
	GetResponse(id int) (models.Response, error)
	SaveRequest(id int) (models.Request, error)
	SaveResponse(id int) (models.Response, error)
}
