package usecase

import (
	"github.com/DmitriyKomarovCoder/http_proxy/internal/models"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api"
)

type Usecase struct {
	repo api.Repository
}

func NewUsecase(repo api.Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) AllRequest() ([]models.Request, error) {
	return []models.Request{}, nil
}

func (u *Usecase) GetRequest(id int) (models.Request, error) {
	return models.Request{}, nil
}

func (u *Usecase) Repeat(id int) (models.Response, error) {
	return models.Response{}, nil
}

func (u *Usecase) Scan(id int) error {
	return nil
}
