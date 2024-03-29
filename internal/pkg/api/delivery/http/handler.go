package http

import (
	"errors"
	"github.com/DmitriyKomarovCoder/http_proxy/common/logger"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"net/http/httputil"
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
	requests, err := h.useCase.AllRequest()

	if errors.Is(err, mongo.ErrNoDocuments) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, "not found documents")
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, requests)
}

func (h *Handler) GetRequest(ctx *gin.Context) {
	id := ctx.Param("id")
	request, err := h.useCase.GetRequest(id)

	if errors.Is(err, mongo.ErrNoDocuments) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, "not found document")
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, request)
}

func (h *Handler) Repeat(ctx *gin.Context) {
	id := ctx.Param("id")
	repeat, err := h.useCase.Repeat(id)

	if errors.Is(err, mongo.ErrNoDocuments) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, "not found document")
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	client := http.Client{}

	resp, err := client.Do(repeat)
	if err != nil {
		h.log.Println("error sending request to target:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)

		return
	}

	var b = []byte{}
	if b, err = httputil.DumpResponse(resp, true); err == nil {
		log.Printf("target response:\n%s\n", string(b))
	}
	ctx.JSON(http.StatusOK, string(b))
}

func (h *Handler) Scan(ctx *gin.Context) {
	id := ctx.Param("id")
	scanInfo, err := h.useCase.Scan(id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, "not found document")
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if scanInfo {
		ctx.JSON(http.StatusOK, "XXE - detected")
		return
	}
	ctx.JSON(http.StatusOK, "XXE not found")
}
