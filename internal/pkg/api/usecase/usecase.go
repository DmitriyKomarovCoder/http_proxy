package usecase

import (
	"github.com/DmitriyKomarovCoder/http_proxy/internal/models"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Usecase struct {
	repo api.Repository
}

func NewUsecase(repo api.Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) AllRequest() ([]models.Request, error) {
	return u.repo.GetAll()
}

func (u *Usecase) GetRequest(id string) (models.Request, error) {
	return u.repo.GetRequest(id)
}

func (u *Usecase) Repeat(id string) (models.Response, error) {
	return models.Response{}, nil
}

func (u *Usecase) Scan(id string) error {
	return nil
}

func (u *Usecase) SaveRequest(request *http.Request, bodyBytes []byte) (primitive.ObjectID, error) {
	postParam, err := utils.ParsePostParams(request)
	if err != nil {
		return primitive.NilObjectID, err
	}

	cleanedBody := utils.CleanNonUTF8(bodyBytes)

	var requestModel = models.Request{
		ID:         primitive.NewObjectID(),
		Method:     request.Method,
		Path:       request.URL.Path,
		Host:       request.Host,
		GetParams:  utils.ParseGetParams(request),
		PostParams: postParam,
		Headers:    utils.ParseHeaders(request.Header),
		Cookie:     utils.ParseCookie(request.Cookies()),
		Body:       string(cleanedBody),
	}

	return u.repo.SaveRequest(requestModel)
}

func (u *Usecase) SaveResponse(id primitive.ObjectID, response http.Response, bodyBytes []byte) error {
	cleanedBody := utils.CleanNonUTF8(bodyBytes)
	var responseModel = models.Response{
		ID:      id,
		Code:    response.StatusCode,
		Headers: utils.ParseHeaders(response.Header),
		Cookie:  utils.ParseCookie(response.Cookies()),
		Body:    string(cleanedBody),
	}

	return u.repo.SaveResponse(responseModel)
}
