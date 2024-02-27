package utils

import (
	"fmt"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/models"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

const (
	HeaderContentTypeApplication = "application/x-www-form-urlencoded"
	HeaderContentType            = "Content-Type"
)

func ParseGetParams(request *http.Request) []models.Params {
	params := make([]models.Params, 0, len(request.URL.Query()))
	queryParams := request.URL.Query()
	for key, values := range queryParams {
		for _, value := range values {
			params = append(params, models.Params{
				Key:   key,
				Value: value,
			})
		}

	}
	return params
}

func ParsePostParams(request *http.Request) ([]models.Params, error) {
	params := make([]models.Params, 0)
	if request.Header.Get(HeaderContentType) == HeaderContentTypeApplication {
		err := request.ParseForm()
		if err != nil {
			return []models.Params{}, fmt.Errorf("error when trying parse params %w", err)
		}

		form := request.PostForm
		for key, values := range form {
			for _, value := range values {
				params = append(params, models.Params{Key: key, Value: value})
			}
		}
	}
	return params, nil
}

func ParseHeaders(headers http.Header) []models.Params {
	masHeaders := make([]models.Params, 0)

	for key, values := range headers {
		for _, value := range values {
			masHeaders = append(masHeaders, models.Params{Key: key, Value: value})
		}
	}

	return masHeaders
}

func ParseCookie(cookies []*http.Cookie) []models.Params {
	masCookies := make([]models.Params, 0, len(cookies))
	for _, v := range cookies {
		masCookies = append(masCookies, models.Params{
			Key:   v.Name,
			Value: v.Value,
		})
	}
	return masCookies
}

func ParseBody(bodyReader io.Reader) (string, error) {
	bodyBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		return "", fmt.Errorf("error when trying parse body %w", err)
	}

	return string(bodyBytes), nil
}

func ChangeRequestToTarget(req *http.Request, targetHost string) {
	targetUrl := addrToUrl(targetHost)
	targetUrl.Path = req.URL.Path
	targetUrl.RawQuery = req.URL.RawQuery
	req.URL = targetUrl

	req.RequestURI = ""
}

func addrToUrl(addr string) *url.URL {
	if !strings.HasPrefix(addr, "https") {
		addr = "https://" + addr
	}
	u, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	return u
}

func CleanNonUTF8(input []byte) []byte {
	var output []byte
	for len(input) > 0 {
		r, size := utf8.DecodeRune(input)
		if r != utf8.RuneError {
			output = append(output, input[:size]...)
		}
		input = input[size:]
	}
	return output
}
