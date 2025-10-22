package controller

import (
	"bme/internal/service"
	"bme/pkg/errorext"
	"bme/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type HeaderEntity struct {
	UserID uint `header:"X-Auth-User-ID"`
}

type HeaderEntityBindingRequired struct {
	UserID uint `header:"X-Auth-User-ID" binding:"required"`
}

type PaginationRequest struct {
	CurrentPage int `form:"current_page"`
	PerPage     int `form:"per_page"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
}

func writeValidationErrors(ctx *gin.Context, validationError *utils.ValidationError) {
	ctx.JSON(http.StatusOK, *validationError)
}

func writeBindingErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
}

func writeBadRequestErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
}

func writeErrorResponse(ctx *gin.Context, err error, logger *logrus.Entry) {
	var extErr *errorext.Error

	if errors.As(err, &extErr) {
		logger.
			WithField("path", ctx.Request.URL.Path).
			WithField("status", extErr.StatusCode).
			WithField("stackTrace", extErr.StackTrace).
			Error(extErr.Error())

		ctx.JSON(extErr.StatusCode, gin.H{"message": extErr.Error(), "key": extErr.Key})

		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
}

func convertQ[T uint](q *string) (*T, *string) {
	if q == nil {
		return nil, nil
	}

	numericVal, err := strconv.Atoi(*q)
	if err != nil {
		return nil, q
	}

	numericValAsUint := T(uint(numericVal))
	return &numericValAsUint, nil
}

func qAs(q *string) (idStartsWith *string, titleStartsWith *string) {
	numericVal, q := convertQ[uint](q)
	if numericVal != nil {
		tmp := strconv.Itoa(int(*numericVal))
		idStartsWith = &tmp
	}

	if q != nil {
		titleStartsWith = q
	}

	return idStartsWith, titleStartsWith
}

func toViewPaginationMetaFromSvc(meta service.PaginationMeta) PaginationMeta {
	return PaginationMeta(meta)
}

func (req PaginationRequest) isEmpty() bool {
	return req.CurrentPage == 0 && req.PerPage == 0
}

func (req PaginationRequest) toSvc() *service.PaginationRequest {
	if req.isEmpty() {
		return service.NewPaginationRequest(service.DefaultCurrentPage, service.DefaultPerPage)
	}

	return &service.PaginationRequest{
		CurrentPage: req.CurrentPage,
		PerPage:     req.PerPage,
	}
}
