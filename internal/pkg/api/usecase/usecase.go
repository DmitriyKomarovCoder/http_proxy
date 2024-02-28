package usecase

import (
	"fmt"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/models"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"strings"
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

func (u *Usecase) Repeat(id string) (*http.Request, error) {
	modelRepeat, err := u.GetRequest(id)
	if err != nil {
		return &http.Request{}, nil
	}

	var body io.Reader
	if modelRepeat.Body != "" {
		body = strings.NewReader(modelRepeat.Body)
	}

	url := fmt.Sprintf("%s://%s%s", modelRepeat.Scheme, modelRepeat.Host, modelRepeat.Path)
	for i, v := range modelRepeat.GetParams {
		if i == 0 {
			url += "?"
		}
		url += v.Key + "=" + v.Value
	}

	newRequest, err := http.NewRequest(modelRepeat.Method, url, body)

	for key, values := range modelRepeat.Headers {
		for _, value := range values {
			masHeaders = append(masHeaders, models.Params{Key: key, Value: value})
		}
	}

	return newRequest, nil
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
		Scheme:     request.URL.Scheme,
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
