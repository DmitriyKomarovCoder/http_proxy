package http

import (
	"github.com/DmitriyKomarovCoder/http_proxy/common/logger"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase api.Usecase
	log     logger.Logger
}

func NewHandler(useCase api.Usecase, log logger.Logger) *Handler {
	return &Handler{
		useCase: useCase,
		log:     log,
	}
}

func (h *Handler) AllRequest(ctx *gin.Context) {

}

func (h *Handler) GetRequest(ctx *gin.Context) {

}

func (h *Handler) Repeat(ctx *gin.Context) {

}

func (h *Handler) Scan(ctx *gin.Context) {

}
