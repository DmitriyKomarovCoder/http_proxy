package proxy

import (
	"github.com/DmitriyKomarovCoder/http_proxy/common/logger"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api"
	"net/http"
)

type ProxyHandler struct {
	apiUsecase api.Usecase
	log        logger.Logger
}

func NewProxy(apiUsecase api.Usecase, log logger.Logger) *ProxyHandler {
	return &ProxyHandler{
		apiUsecase: apiUsecase,
		log:        log,
	}
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	return
}
