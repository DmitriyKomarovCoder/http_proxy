package repository

import (
	"github.com/DmitriyKomarovCoder/http_proxy/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetAll() ([]models.Request, error) {
	return []models.Request{}, nil
}

func (r *Repo) GetRequest(id int) (models.Request, error) {
	return models.Request{}, nil
}

func (r *Repo) GetResponse(id int) (models.Response, error) {
	return models.Response{}, nil
}

func (r *Repo) SaveRequest(id int) (models.Request, error) {
	return models.Request{}, nil
}

func (r *Repo) SaveResponse(id int) (models.Response, error) {
	return models.Response{}, nil
}
